#

```shell
docker run -it busybox /bin/sh
# -it 分配tty
# 在容器中执行 /bin/sh

# 运行ps
# PID  USER   TIME COMMAND
# 1 root   0:00 /bin/sh
# 10 root   0:00 ps
# 可以发现/bin/sh的pid为1
```

原理,
linux中系统创建线程

```c
int pid = clone(main_function,stack_size,SIGCHLD,NULL);

// 可以在参数中指定CLONE_NEWPID
int pid = clone(main_function,stack_size,CLONE_NEWPID | SIGCHLD,NULL);
// 这时,新创建的这个进程将会"看到"一个全新的进程空间
// 多次执行这个clone()时, 就会创建多个PID Namespace
// 每个Namespace中的应用进程, 既看不到宿主机里真正的进程空间, 也看不到其他PID Namespace里的具体情况
// 除了PID Namespace, Linux还有UTS,IPC,Network,User这些Namespace
```

在理解Namespace的工作方式之后, 跟真实存在的虚拟机不同  
在使用docker的时候, 并没有真正的"Docker 容器"运行在宿主机里  
Docker项目帮助用户启动的, 还是原来的应用进程, 只不过在创建这些进程时, Docker为它们加上了各种各样的Namespace参数  

## docker隔离与限制  

container控制组  
CPU子系统
    -/sys/fs/cgroup/cpu/container/cpu.cfs_quota_us <- 默认-1,无限制
    -/sys/fs/cgroup/cpu/container/cpu.cfs_period_us <- 默认100ms(100000us)
    -/sys/fs/cgroup/cpu/container/tasks <- 被限制进程的PID写入这个文件,会立即生效

除CPU子系统外, Cgroups还有
    - blkio -> 为块设备设定I/O限制, 一般用于磁盘等设备
    - cpuset -> 为进程分配单独的CPU核和对应的内存节点
    - memory -> 为进程设定内存使用的限制

在执行`docker run`时可以添加参数,对一些系统资源进行限制  
如:

```shell
docker run -it --cpu-period=100000 --cpu-quota=20000 ubuntu /bin/bash
```

对文件的隔离  

创建新进程时,除了声明要启用Mount Namespace之外, 还需要告诉容器进程哪些目录需要重新挂载
如:  

```shell
mount("","/",NULL,MS_PRIVATE,"");
mount("none","/tmp","tmpfs",0,"");
```

这就是Mount Namespace跟其他Namespace的使用略有不同的地方, 它对容器进程视图的改变, 一定是伴随着挂载操作(mount)才能生效.(Linux中的chroot)
这个挂载在容器根目录上, 用来为容器进程提供合理后执行环境的文件系统就是所谓的"容器镜像",叫做 -> rootfs(根文件系统)  

Docker在镜像的设计中, 引入了层(layer)的概念. 也就是说, 用户制作镜像的每一步操作, 都会生成一个层, 也就是一个增量rootfs.  
用到了一种叫作联合文件系统(Union File System)的能力, 也叫UnionFS, 最主要的功能是将多个不同位置的目录联合挂载(union mount)到同一个目录下.  
比如使用联合挂载的方式, 将这两个目录挂载到一个公共的目录C上

```shell
# .
# ├── A
# │ ├── a
# │ └── x
# └── B
#   ├── b
#   └── x

mount -t aufs -o dirs=./A:./B none ./C

# ./C
# ├── a
# ├── b
# └── x

# A和B中的x被合并为一个文件夹c
```

Docker默认使用的是AuFS这个联合文件系统的实现.(全程Another UnionFS 后改名Alternative UnionFS)

例子:  
Ubuntu镜像由五个层构成  

1. 只读层 <- 只读
2. Init层
    - 夹在只读层和读写层之间
    - 专门用来存放/etc/hosts, /etc/resolv.conf
3. 可读写层
    - 删只读层里的一个文件时, AuFS会在可读写层创建一个whiteout文件, 把只读层里的文件"遮挡"起来
    - 所以这个可读写的作用就是专门存放你修改rootfs后产生的增量

## 制作容器镜像

### Dockerfile

```dockerfile
# 使用官方提供的 Python 开发镜像作为基础镜像
FROM python:2.7-slim
# 将工作目录切换为 /app
WORKDIR /app
# 将当前目录下的所有内容复制到 /app 下
ADD . /app
# 使用 pip 命令安装这个应用所需要的依赖
RUN pip install --trusted-host pypi.python.org -r requirements.txt
# 允许外界访问容器的 80 端口
EXPOSE 80
# 设置环境变量
ENV NAME World
# 设置容器进程为:python app.py,即:这个 Python 应用的启动命令
CMD ["python", "app.py"]
# 另外,在使用 Dockerfile 时,你可能还会看到一个叫作 ENTRYPOINT 的原
# 语。实际上,它和 CMD 都是 Docker 容器进程启动所必需的参数,完整执行
# 格式是:“ENTRYPOINT CMD”。
# 但是,默认情况下,Docker 会为你提供一个隐含的 ENTRYPOINT,即:
# /bin/sh -c。所以,在不指定 ENTRYPOINT 时,比如在我们这个例子里,实际
# 上运行在容器里的完整进程是:/bin/sh -c “python app.py”,即 CMD 的内容就
# 是 ENTRYPOINT 的参数。
```

在这样一个目录中

```shell
Dockerfile app.py requirements.txt
````

制作Docker镜像

```shell
docker build -t helloworld .

# -t 的作用是给这个镜像加一个Tag,相当于这个镜像的名字
```

使用一个镜像, 如上一步制作的helloworld

```shell
docker run -p 4000:80 helloworld
# 如果Dockerfile中没有指定CMD, 则需要
docker run -p 4000:80 helloworld python app.py
# -p 4000:80 是将docker容器内的80端口映射在宿主机的4000端口上
# 或者使用dokcer inspect查看容器的IP地址, 然后访问 http://<容器IP地址>:80
```

上传镜像到Docker Hub

```shell
# 注册一个Docker Hub帐号 使用docker login 命令登陆
# 用docker tag命名镜像
docker tag helloworld <上面注册的用户名>/helloworld:v1
# 上传
docker push <用户名>/helloworld:v1
```

使用docker commit指令, 把一个正在运行的容器, 直接体检为一个镜像, 比如

```shell
docker exec -it 4ddf4638572d /bin/sh
# 在容器内新建一个文件
root@4ddf4638572d:/app# touch test.txt
root@4ddf4638572d:/app# exit
# 将这个新建的文件提交到镜像中保存
docker commit 4ddf4638572d <用户名>/helloworld:v2
```

docker专门提供了一个参数可以 启动一个容器并"加入"到另一个容器的Network Namespace中, 这个参数就是 -net

```shell
docker run -it --net container:4ddf4638572d busybox ifconfig
# 如果指定 -net=host 就意味着这个容器不会为进程启用 Network Namespace, 这就意味着,这个容器拆除了 Network Namespace 的“隔离墙”,所以,它会和宿主机上的其他普通进程一样,直接共享宿主机的网络栈。这就为容器直接操作和使用宿主机网络提供了一个渠道。
```

Docker Volume机制, 允许你将宿主机上指定的目录或者文件, 挂载到容器里面进行读取和修改操作

```shell
docker run -v /test ...
docker run -v /home:/test ...
# 第一种没有显式声明宿主机目录, 那么Docker就会默认在宿主机上创建一个临时目录 /var/lib/docker/volumes/[VOLUME_ID]/_data, 然后把它挂载到容器的 /test 目录上
# 第二种Docker就直接把宿主机的 /home 目录挂载到容器的 /test 目录上
```

Volume用了上面提到的联合挂载机制进行实现

## kubernetes

架构

```markdown
+----+
|ETCD| <- (用来存储的)
+----+
   |
   |     grpc
   +--------------+
                  |
+-------------------------------------+
|       Master    |                   |
|  +----------+ +------+ +----------+ |
|  |Controller| | API  | | Scheduler| |
|  | Manager  | |Server| |          | |
|  +----------+ +------+ +----------+ |
|                  |                  |
+-------------------------------------+
                   |
                   +--------+
                            | protobuf
    +--------------+        |
    |    Node      |--------+
    | +----------+ |
    | |Networking| |
    | +----------+ |
    | +--------+   |
    | | kubelet|   |
    | +--------+   |
    | +---------+  |   ...... <- 很多node节点
    | |Container|  |
    | | Runtime |  |
    | +---------+  |
    | +------+     |
    | |Volume|     |
    | |Plugin|     |
    | +------+     |
    | +------+     |
    | |Device|     |
    | |Plugin|     |
    | +------+     |
    | +--------+   |
    | |Linux OS|   |
    | +--------+   |
    +--------------+
```

比如:一个“Web 容器”和它要访问的数据库“DB 容器”。  
在常规环境下,这些应用往往会被直接部署在同一台机器上,通过 Localhost通信,通过本地磁盘目录交换文件。而在 Kubernetes 项目中,这些容器则会被划分为一个“Pod”,Pod 里的容器共享同一个 Network Namespace、同一组数据卷,从而达到高效率交换信息的目的。Pod是Kubernetes项目中最基础的一个对象.  

而对于另外一种更为常见的需求,比如 Web 应用与数据库之间的访问关系,Kubernetes 项目则提供了一种叫作“Service”的服务。  Kubernetes 项目的做法是给 Pod 绑定一个 Service 服务,而 Service服务声明的 IP 地址等信息是“终生不变”的。这个 Service 服务的主要作用,就是作为 Pod 的代理入口(Portal),从而代替 Pod 对外暴露一个固定的网络地址。  

Kubernetes启动一个容器化任务

编写一个yaml文件,比如

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: nginx-deployment
    labels:
        app: nginx
spec:
    replicas: 2
    selector:
        matchLabels:
            app: nginx
    template:
        metadata:
            labels:
                app: nginx
        spec:
            containers:
            - name: nginx
                image: nginx:1.7.9
                ports:
                - containerPort: 80
```

在上面这个 YAML 文件中,我们定义了一个 Deployment 对象,它的主体部
分(spec.template 部分)是一个使用 Nginx 镜像的 Pod,而这个 Pod 的副
本数是 2(replicas=2)。  
然后执行

```shell
kubectl create -f nginx-deployment.yaml
```

kubernetes一键部署工具kubeadm

```shell
# 创建一个Mater节点
kubeadm init
# 将一个Node节点加入到当前的集群中
kubeadm join <Master 节点的IP和端口>
```

kubeadm工作原理 -> 把 kubelet 直接运行在宿主机上,然后使用容器部署其他的 Kubernetes 组件。

kubeadm工作流程

- 一系列检查工作,以确定这台机器可以用来部署Kubernetes,称为"Preflight Checks"
- 生成Kubernetes对外提供服务所需的各种证书和对应的目录  
      Kubernetes对外提供服务时,除非专门开启"不安全模式",否则都要通过HTTPS才能访问kube-apiserver
- kubeadm接下来会为其他组件生成访问kube-apiserver所需的配置文件
- 接下来,kubeadm会为Master组件生成Pod配置文件
- Master 容器启动后,kubeadm 会通过检查 localhost:6443/healthz 这个Master 组件的健康检查 URL,等待 Master 组件完全运行起来
- 然后,kubeadm 就会为集群生成一个 bootstrap token。在后面,只要持有这个 token,任何一个安装了 kubelet 和 kubadm 的节点,都可以通过kubeadm join 加入到这个集群当中。

使用`kubeadm init`部署时使用自定义启动参数

```shell
kubeadm init --config kubeadm.yaml
```

最后,我再来回答一下我在今天这次分享开始提到的问题:kubeadm能够用于生产环境吗?  
到目前为止(2018 年 9 月),这个问题的答案是:不能。

通过Talint/Toleration调整Master执行Pod的策略

它的原理非常简单:一旦某个节点被加上了一个 Taint,即被“打上了污点”,那
么所有 Pod 就都不能在这个节点上运行,因为 Kubernetes 的 Pod 都有“洁
癖”。  
除非,有个别的 Pod 声明自己能“容忍”这个“污点”,即声明了 Toleration,它
才可以在这个节点上运行。  
为节点打上"Taint"

```shell
kubectl taint nodes node1 foo=bar:NoSchedule
# 这时,该 node1 节点上就会增加一个键值对格式的 Taint,即:
# foo=bar:NoSchedule。其中值里面的 NoSchedule,意味着这个 Taint 只会在调度新 Pod 时产生作用,而不会影响已经在 node1 上运行的 Pod,哪怕它们没有 Toleration。

```

删除Taint

```shell
kubectl taint nodes --all node-role.kubernetes.io/master-
```

kubernetes yaml文件解析

```yaml
apiVersion: apps/v1
# kind 指定了这个API对象的类型,是一个Deployment
# 所谓Deployment,是一个定义多副本应用的对象
kind: Deployment
metadata:
    name: nginx-deployment
spec:
    selector:
        matchLabels:
            app: nginx
        # spec.replicas 定义副本的个数,这里为2
        replicas: 2
        # spec.template定义一个Pod的模板,这个模板描述了我想要创建的Pod的细节
        template:
            # metadata <- API对象的"标识"
            metadata:
                # labels 比如这个, 会把所有正在运行的, 携带"app:nginx"标签的Pod识别为被管理的对象,并确保这些Pod的总数严格等于两个
                labels:
                    app: nginx
                # 另外,在 Metadata 中,还有一个与 Labels 格式、层级完全相同的字段叫Annotations,它专门用来携带 key-value 格式的内部信息。所谓内部信息,指的是对这些信息感兴趣的,是 Kubernetes 组件本身,而不是用户。所以大多数 Annotations,都是在 Kubernetes 运行过程中,被自动加在这个 API 对象上。
            spec:
                containers:
                - name: nginx
                    image: nginx:1.7.9
                    ports:
                    - containerPort: 80
```

运行这个配置文件

```shell
kubectl create -f nginx-deployment.yaml
# 通过kubectl get检查这个YAML文件运行起来的状态是不是与我们预期的一致
kubectl get pods -l app=nginx
# 使用kubectl describe 查看一个API对象的细节
kubectl describe pod nginx-deployment-67594d6bf6-9gdvr
```

对服务进行升级

- 修改YAML文件

    ```yaml
    ...
        spec:
            containers:
            - name: nginx
                image: nginx:1.8 # 这里被从 1.7.9 修改为 1.8
                ports:
            - containerPort: 80
    ```

- 然后用kubectl replace完成更新

    ```shell
    kubectl replace -f nginx-deployment.yaml
    ```

- 或者使用kubectl apply,统一进行kubernetes对象的创建和更新操作

    ```shell
    kubectl apply -f nginx-deployment.yaml
    # 修改nginx-deployment.yaml的内容
    kubectl apply -f nginx-deployment.yaml
    ```

使用kubectl exec指令,进入到这个Pod当中

```shell
kubectl exec -it nginx-deployment-5c678cfb6d-lg9lw -- /bin/bash
```

从kubernetes集群中删除Nginx Deployment

```shell
kubectl delete -f nginx-deployment.yaml
```

Pod的实现原理

- 关于Pod最重要的一个事实是它只是一个逻辑概念  
  kubernetes真正处理的,还是宿主机操作系统上Linux容器的Namespace和Cgroup,而并不存在一个所谓的Pod边界或者隔离环境
- Pod,其实是一组共享了某些资源的容器  
  Pod里的所有容器,共享的是同一个Network Namespace, 并且可以声明共享同一个Volume
- 在Kubernetes项目里Pod的实现需要使用一个中间容器, 这个容器叫做Infra容器. 在这个Pod中,Infra永远都是第一个被创建的容器,而其他用户定义的容器,则通过Join Network Namespace的方式, 与Infra容器关联在一起.
- 在 Kubernetes 项目里,Infra 容器一定要占用极少的资源,所以它使用的是一个非常特殊的镜像,叫作:k8s.gcr.io/pause。这个镜像是一个用汇编语言编写的、永远处于“暂停”状态的容器,解压后的大小也只有100~200 KB 左右。
- Pod的声明周期只跟Infra容器一致, 而与容器A和B无关

```yaml
apiVersion: v1
kind: Pod
...
spec:
    nodeSelector:
        disktype: ssd
# 这样的一个配置 意味着这个Pod永远只能运行在携带了"disktype: ssd"标签的节点上;否则, 他将调度失败.

apiVersion: v1
kind: Pod
...
spec:
    hostAliases:
    - ip: "10.1.2.3"
hostnames:
    - "foo.remote"
    - "bar.remote"
...
# HostAliases 定义了Pod的hosts文件(比如/etc/hosts)里的内容
# cat /etc/hosts
# Kubernetes-managed hosts file.
# 127.0.0.1 localhost
# ...
# 10.244.135.10 hostaliases-pod
# 10.1.2.3 foo.remote
# 10.1.2.3 bar.remote

apiVersion: v1
kind: Pod
metadata:
    name: nginx
spec:
    shareProcessNamespace: true
    containers:
    - name: nginx
        image: nginx
    - name: shell
        image: busybox
        stdin: true
        tty: true
# 设置了shareProcessNamespace=true, 意味着这个Pod里的容器要共享PID Namespace
# tty,stdin 就是docker run -it 中的 -it

apiVersion: v1
kind: Pod
metadata:
    name: nginx
spec:
    hostNetwork: true
    hostIPC: true
    hostPID: true
containers:
- name: nginx
    image: nginx
- name: shell
    image: busybox
    stdin: true
    tty: true
# 这个Pod的 定义共享主机的Network,IPC,PID Namespace, 这就意味着, 这个Pod里的所有容器, 会直接使用宿主机的网络, 直接与宿主机进行IPC通信, 看到宿主机里正在运行的所有进程.

apiVersion: v1
kind: Pod
metadata:
    name: lifecycle-demo
spec:
    containers:
    - name: lifecycle-demo-container
        image: nginx
        lifecycle:
            postStart:
                exec:
                    command: ["/bin/sh", "-c", "echo Hello from the postStart handler > /usr/share/message"]
            preStop:
                exec:
                    command: ["/usr/sbin/nginx","-s","quit"]
# Lifecycle字段, 定义的是Container Lifecycle Hooks, 顾名思义就是在容器状态发生变化时触发一系列"钩子"
# postStart 在容器启动后, 立刻执行一个指定的操作
# preStop 在容器被杀死之前, 执行操作
#
# ImagePullPolicy字段 定义了拉取的策略
# 值                 作用
# Always        即每次创建Pod都重新拉取一次镜像
# Never         意味着Pod永远不会主动拉取这个镜像
# ifNotPresent  只在宿主机上不存在这个镜像时才拉取
```

Pod对象在Kubernetes中的生命周期  
Pod 生命周期的变化,主要体现在 Pod API 对象的 Status 部分,这是它除了Metadata 和 Spec 之外的第三个重要字段。其中,pod.status.phase,就是Pod 的当前状态,它有如下几种可能的情况:

- Pending。这个状态意味着,Pod 的 YAML 文件已经提交给了Kubernetes,API 对象已经被创建并保存在 Etcd 当中。但是,这个Pod 里有些容器因为某种原因而不能被顺利创建。比如,调度不成功。
- Running。这个状态下,Pod 已经调度成功,跟一个具体的节点绑定。它包含的容器都已经创建成功,并且至少有一个正在运行中。
- Succeeded。这个状态意味着,Pod 里的所有容器都正常运行完毕,并且已经退出了。这种情况在运行一次性任务时最为常见。
- Failed。这个状态下,Pod 里至少有一个容器以不正常的状态(非 0 的返回码)退出。这个状态的出现,意味着你得想办法 Debug 这个容器的应用,比如查看 Pod 的 Events 和日志。
- Unknown。这是一个异常状态,意味着 Pod 的状态不能持续地被kubelet 汇报给 kube-apiserver,这很有可能是主从节点(Master 和Kubelet)间的通信出现了问题。

Projected Volume  

在 Kubernetes 中,有几种特殊的 Volume,它们存在的意义不是为了存放容器里的据,也不是用来进行容器和宿主机之间的数据交换。这些特殊Volume 的作用,是为容器提供预先定义好的数据。所以,从容器的角度来看,这些 Volume 里的信息就是仿佛是被 Kubernetes“投射”(Project)进入容器当中的。这正是 Projected Volume 的含义。

第一种 -> **Secret**

```yaml
apiVersion: v1
kind: Pod
metadata:
    name: test-projected-volume
spec:
    containers:
    - name: test-secret-volume
        image: busybox
        args:
        - sleep
        - "86400"
        volumeMounts:
        - name: mysql-cred
            mountPath: "/projected-volume"
            readOnly: true
        volumes:
        - name: mysql-cred
            projected:
                sources:
                - secret:
                    name: user
                - secret:
                    name: pass
# 在这个 Pod 中,我定义了一个简单的容器。它声明挂载的 Volume,并不是常见的 emptyDir 或者 hostPath 类型,而是 projected 类型。而这个Volume 的数据来源(sources),则是名为 user 和 pass 的 Secret 对象,分别对应的是数据库的用户名和密码。
```

创建Secret对象

```shell
kubectl create secret generic user --from-file=./username.txt
kubectl create secret generic user --from-file=./password.txt
# 查看secret对象
kubectl get secrets
# 通过直接编写YAML文件的方式创建这个Secret对象
# apiVersion: v1
# kind: Secret
# metadata:
    # name: mysecret
# type: Opaque
# data:
    # user: YWRtaW4=
    # pass: MWYyZDFlMmU2N2Rm
```

像这样通过挂载方式进入到容器里的 Secret,一旦其对应的Etcd 里的数据被更新,这些 Volume 里的文件内容,同样也会被更新。其实,这是 kubelet 组件在定时维护这些 Volume。  

**ConfigMAP**  

与 Secret 类似的是 ConfigMap,它与 Secret 的区别在于,ConfigMap 保存的是不需要加密的、应用所需的配置信息。而 ConfigMap 的用法几乎与Secret 完全相同:你可以使用 kubectl create configmap 从文件或者目录创建ConfigMap,也可以直接编写 ConfigMap 对象的 YAML 文件。

比如,一个 Java 应用所需的配置文件(.properties 文件),就可以通过下面
这样的方式保存在 ConfigMap 里:

```shell
# .properties 文件的内容
$ cat example/ui.properties
color.good=purple
color.bad=yellow
allow.textmode=true
how.nice.to.look=fairlyNice
# 从.properties 文件创建 ConfigMap
$ kubectl create configmap ui-config --from-file=example/ui.properties
# 查看这个 ConfigMap 里保存的信息 (data)
$ kubectl get configmaps ui-config -o yaml
apiVersion: v1
data:
    ui.properties: |
        color.good=purple
        color.bad=yellow
        allow.textmode=true
        how.nice.to.look=fairlyNice
kind: ConfigMap
metadata:
    name: ui-config
...
```

**Downward API** -> 让 Pod 里的容器能够直接获取到这个 Pod API 对象本身的信息

```yaml
apiVersion: v1
kind: Pod
metadata:
    name: test-downwardapi-volume
    labels:
        zone: us-est-coast
        cluster: test-cluster1
        rack: rack-22
spec:
    containers:
        - name: client-container
            image: k8s.gcr.io/busybox
            command: ["sh", "-c"]
            args:
            - while true; do
                if [[ -e /etc/podinfo/labels ]]; then
                    echo -en '\n\n'; cat /etc/podinfo/labels; fi;
                sleep 5;
            done;
        volumeMounts:
            - name: podinfo
                mountPath: /etc/podinfo
                readOnly: false
        volumes:
            - name: podinfo
            projected:
                sources:
                - downwardAPI:
                    items:
                        - path: "labels"
                            fieldRef:
                                fieldPath: metadata.labels
# 在这个 Pod 的 YAML 文件中,我定义了一个简单的容器,声明了一个projected 类型的 Volume。只不过这次 Volume 的数据来源,变成了Downward API。而这个 Downward API Volume,则声明了要暴露 Pod 的metadata.labels 信息给容器。
# 通过这样的声明方式,当前 Pod 的 Labels 字段的值,就会被 Kubernetes自动挂载成为容器里的 /etc/podinfo/labels 文件。
# Downward API还支持更多暴露其他字段, 可以去官方文档查看
# 需要注意的是,Downward API 能够获取到的信息,一定是 Pod 里的容器进程启动之前就能够确定下来的信息
```

**Service Account**  
Service Account 对象的作用,就是 Kubernetes 系统内置的一种“服务账户”,它是 Kubernetes 进行权限分配的对象。比如,Service Account A,可以只被允许对 Kubernetes API 进行 GET 操作,而 Service Account B,则可以有Kubernetes API 的所有操作的权限。  
像这样的 Service Account 的授权信息和文件,实际上保存在它所绑定的一个特殊的 Secret 对象里的。这个特殊的 Secret 对象,就叫作ServiceAccountToken。  
所以说,Kubernetes 项目的 Projected Volume 其实只有三种,因为第四种ServiceAccountToken,只是一种特殊的 Secret 而已。  

实现原理

```shell
#查看一下任意运行在Kubernetes集群的Pod
kubectl describe pod nginx-deployment-5c678cfb6d-lg91w
Containers:
...
    Mounts:
        /var/run/secrets/kubernetes.io/serviceaccount from default-token-s8rbq (ro)
    Volumes:
        default-token-s8rbq:
        Type:   Secret (a volume populated by a Secret)
        SecretName: default-token-s8rbq
        Optional: false
```

就会发现每个Pod都已经自动声明一个类型是Secret, 名为default-token-xxxx的Volume, 然后自动挂载在每个容器的一个固定目录上.  
这个Secret类型的Volume, 正是默认 Service Account 对应的ServiceAccountToken, 所以 Kubernetes其实在每个Pod创建的时候, 自动在它的spec.volumes部分添加上了默认ServiceAccountToken的定义, 然后自动给每个容器加上了对应的volumeMounts字段, 这个过程对于用户来说完全是透明的.  
这样,一旦 Pod 创建完成,容器里的应用就可以直接从这个默认ServiceAccountToken 的挂载目录里访问到授权信息和文件。这个容器内的路径在 Kubernetes 里是固定的,即:`/var/run/secrets/kubernetes.io/serviceaccount`  
这种把 Kubernetes 客户端以容器的方式运行在集群里,然后使用 defaultService Account 自动授权的方式,被称作“InClusterConfig” <- 原文章作者最推荐进行Kubernetes API的授权方式

**容器健康检查和恢复机制**  

```yaml
apiVersion: v1
kind: Pod
metadata:
    labels:
        test: liveness
    name: test-liveness-exec
spec:
    containers:
    - name: liveness
        image: busybox
        args:
        - /bin/sh
        - -c
        - touch /tmp/healthy; sleep 30; rm -rf /tmp/healthy; sleep 600
        livenessProbe: # 健康检查
            exec:
                command:
                - cat
                - /tmp/healthy
            initialDelaySeconds: 5 # 在容器启动5s后开始执行
            periodSeconds: 5 # 每5s执行一次

# 在这个 Pod 中,我们定义了一个有趣的容器。它在启动之后做的第一件事,就是在 /tmp 目录下创建了一个 healthy 文件,以此作为自己已经正常运行的标志。而 30 s 过后,它会把这个文件删除掉。
# 与此同时,我们定义了一个这样的 livenessProbe(健康检查)。它的类型是exec,这意味着,它会在容器启动后,在容器里面执行一句我们指定的命令,比如:“cat /tmp/healthy”。
# 是 Kubernetes 里的 Pod 恢复机制,也叫 restartPolicy。它是Pod 的 Spec 部分的一个标准字段(pod.spec.restartPolicy),默认值是Always,即:任何时候这个容器发生了异常,它一定会被重新创建。 <- 这里注意是重新创建,而不是重启
# Always:在任何情况下,只要容器不在运行状态,就自动重启容器;
# OnFailure: 只在容器 异常时才自动重启容器;
# Never: 从来不重启容器。
```

**kube-controller-manager**  

实际上,这个组件,就是一系列控制器的集合。我们可以查看一下Kubernetes 项目的 pkg/controller 目录:

```shell
$ cd kubernetes/pkg/controller/
$ ls -d */
deployment/ job/ podautoscaler/
cloud/ disruption/ namespace/
replicaset/ serviceaccount/ volume/
cronjob/ garbagecollector/ nodelifecycle/
replication/ statefulset/ daemon/
...
```

这个目录下面的每一个控制器,都以独有的方式负责某种编排功能。而我们的
Deployment,正是这些控制器中的一种。

**Deployment**  

Deployment 看似简单,但实际上,它实现了 Kubernetes 项目中一个非常重要的功能:Pod 的“水平扩展 / 收缩”(horizontal scaling out/in)。这个功能,是从 PaaS 时代开始,一个平台级项目就必须具备的编排能力。  
举个例子,如果你更新了 Deployment 的 Pod 模板(比如,修改了容器的镜像),那么 Deployment 就需要遵循一种叫作“滚动更新”(rolling update)的方式,来升级现有的容器。
而这个能力的实现,依赖的是 Kubernetes 项目中的一个非常重要的概念(API 对象):ReplicaSet。

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
    name: nginx-set
    labels:
        app: nginx
    spec:
        replicas: 3
        selector:
            matchLabels:
            app: nginx
        template:
            metadata:
                labels:
                    app: nginx
            spec:
                containers:
                - name: nginx
                image: nginx:1.7.9
# 一个ReplicaSet对象 其实就是由副本数目的定义和一个Pod模板组成的.
# 更重要的是 Deployment控制器实际操纵的 正是这样的ReplicaSet对象而不是Pod对象
```

Deployment与ReplicaSet以及Pod的关系

```markdown
+------------+
| Deployment |
+------------+
      |
+------------+
| ReplicaSet |  <- 通过控制器模式保证系统中的Pod的个数永远等于指定的个数
+------------+
  /    |   \
 pod  pod  pod
```
