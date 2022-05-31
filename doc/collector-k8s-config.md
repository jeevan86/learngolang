### env

```shell
export KUBE_CLUSTER_ID=kubernetes
export KUBE_API_SERVER=https://192.168.7.154:6443
export KUBE_API_CA_KEY=/etc/kubernetes/pki/ca.key
export KUBE_API_CA_CRT=/etc/kubernetes/pki/ca.crt
export CLIENT_USR_NAME=monitor
export USER_NAME_SPACE=monitor
```

### monitor-rbac.yaml
apiGroups中"" 标明 core API 组
```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: monitor
  labels:
    name: monitor
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: monitor
  name: monitor
rules:
- apiGroups:
    - ""
    - extensions
    - apps
    - batch
  resources:
    - nodes
    - services
    - pods
    - deployments
    - daemonsets
    - replicasets
    - replicationcontrollers
    - statefulsets
    - jobs
    - cronjobs
  verbs:
    - get
    - watch
    - list
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: monitor
  name: read-resources
subjects:
- kind: User
  name: monitor
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: monitor
  apiGroup: rbac.authorization.k8s.io
```

### mk-kube-conf.sh

```shell
export KUBE_CONFIG_FILE=./kube-${KUBE_CLUSTER_ID}-${CLIENT_USR_NAME}.conf

openssl genrsa -out ./${CLIENT_USR_NAME}.key 2048
openssl req -new                                     \
 -subj "/CN=${CLIENT_USR_NAME}/O=${USER_NAME_SPACE}" \
 -key ./${CLIENT_USR_NAME}.key                       \
 -out ./${CLIENT_USR_NAME}.csr
openssl x509 -req -CAcreateserial -days 9999 \
  -CA ${KUBE_API_CA_CRT}                     \
  -CAkey ${KUBE_API_CA_KEY}                  \
 -in ./${CLIENT_USR_NAME}.csr -out ./${CLIENT_USR_NAME}.crt

# cat /etc/kubernetes/pki/ca.crt | base64 -w 0
#   LS0tLS1CRUdJTiBD ...
# cat monitor.crt | base64 -w 0
#   LS0tLS1CRUdJTiBD ...
# cat monitor.key | base64 -w 0
#   LS0tLS1CRUdJTiBS ...

# ============ 客户端配置
# 集群
# https://www.kubernetes.org.cn/doc-52
## --certificate-authority=/etc/kubernetes/pki/ca.crt # 设置kuebconfig配置文件中集群选项中的certificate-authority路径
## --embed-certs=true                                 # 设置kuebconfig配置文件中集群选项中的embed-certs开关
## --insecure-skip-tls-verify=false                   # 设置kuebconfig配置文件中集群选项中的insecure-skip-tls-verify开关
## --server="https://1.2.3.4:6443"                    # 设置kuebconfig配置文件中集群选项中的server
kubectl --kubeconfig=${KUBE_CONFIG_FILE} \
  config set-cluster ${KUBE_CLUSTER_ID}  \
    --server=${KUBE_API_SERVER}          \
    --embed-certs=true                   \
    --certificate-authority=${KUBE_API_CA_CRT} # --insecure-skip-tls-verify=true
    
# 用户
# https://www.kubernetes.org.cn/doc-54
## --client-certificate="": 设置kuebconfig配置文件中用户选项中的证书文件路径。
## --client-key="": 设置kuebconfig配置文件中用户选项中的证书密钥路径。
## --embed-certs=false: 设置kuebconfig配置文件中用户选项中的embed-certs开关。
## --password="": 设置kuebconfig配置文件中用户选项中的密码。
## --token="": 设置kuebconfig配置文件中用户选项中的令牌。
## --username="": 设置kuebconfig配置文件中用户选项中的用户名。
kubectl --kubeconfig=${KUBE_CONFIG_FILE}          \
  config set-credentials ${CLIENT_USR_NAME}       \
    --client-key=./${CLIENT_USR_NAME}.key         \
    --client-certificate=./${CLIENT_USR_NAME}.crt \
    --embed-certs=true                    # --username=admin --password=uXFGweU9l35qcif # --token=bearer_token

# 上下文
# https://www.kubernetes.org.cn/doc-53
## --cluster=""   # 设置kuebconfig配置文件中环境选项中的集群
## --user=""      # 设置kuebconfig配置文件中环境选项中的用户
## --namespace="" # 设置kuebconfig配置文件中环境选项中的命名空间
kubectl --kubeconfig=${KUBE_CONFIG_FILE}                   \
  config set-context ${CLIENT_USR_NAME}@${KUBE_CLUSTER_ID} \
    --cluster=${KUBE_CLUSTER_ID}                           \
    --user=${CLIENT_USR_NAME}                              \
    --namespace=${USER_NAME_SPACE}
    
kubectl --kubeconfig=${KUBE_CONFIG_FILE}                   \
  config set current-context ${CLIENT_USR_NAME}@${KUBE_CLUSTER_ID}
```

### kube-monitor.conf

```yaml
apiVersion: v1
kind: Config
preferences: {}
current-context: monitor@kubernetes
clusters:
  - cluster:
      # certificate-authority: /etc/kubernetes/pki/ca.crt
      # cat ca.crt | base64 -w 0
      certificate-authority-data: LS0tLS1CRUdJTiB ...
      server: https://192.168.1.1:6443
    name: kubernetes
users:
  - name: monitor
    user:
      # client-certificate: /tmp/monitor.crt
      # client-key: /tmp/monitor.key
      # cat monitor.crt | base64 -w 0
      client-certificate-data: LS0tLS1CRUdJTiB ...
      # cat monitor.key | base64 -w 0
      client-key-data: LS0tLS1CRUdJTiBSU0E ...
contexts:
  - name: monitor@kubernetes
    context:
      cluster: kubernetes
      user: monitor
```

### kubectl

```shell
export KUBE_CONFIG_FILE=./kube-${KUBE_CLUSTER_ID}-${CLIENT_USR_NAME}.conf
kubectl --kubeconfig=${KUBE_CONFIG_FILE} \
 get nodes,services,pods,deployments,daemonsets,replicasets,replicationcontrollers,statefulsets,jobs,cronjobs -n runner
NAME                                                   READY   STATUS    RESTARTS   AGE
pod/dind-848d7545c9-c6s52                              1/1     Running   0          102d
pod/dind-848d7545c9-m7jcl                              1/1     Running   0          67d
pod/gitlab-runner-205-gitlab-runner-7b7b57d7d5-j7rhv   1/1     Running   0          76d
pod/gitlab-runner-205-gitlab-runner-7b7b57d7d5-l9fjq   1/1     Running   0          76d
pod/proxy-7q7dj                                        1/1     Running   0          49d

NAME                                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/dind                              2/2     2            2           102d
deployment.apps/gitlab-runner-205-gitlab-runner   2/2     2            2           107d

NAME                   DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR       AGE
daemonset.apps/proxy   1         1         1       1            1           runner-node=proxy   108d

NAME                                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/dind-55b78cd4                                0         0         0       102d
replicaset.apps/dind-5798fd8b9b                              0         0         0       102d
replicaset.apps/dind-845d6f5d7d                              0         0         0       102d
replicaset.apps/dind-848d7545c9                              2         2         2       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-54b988f679   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-57b645f5f9   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-5c997d9796   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-5d5bcfcc58   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-5f6bf89698   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-6468445c5d   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-6f54fbf8fb   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-7b7b57d7d5   2         2         2       76d
replicaset.apps/gitlab-runner-205-gitlab-runner-7bb754f54d   0         0         0       102d
replicaset.apps/gitlab-runner-205-gitlab-runner-7c7f6b6bb8   0         0         0       76d
replicaset.apps/gitlab-runner-205-gitlab-runner-f486b95c     0         0         0       102d

kubectl --kubeconfig=${KUBE_CONFIG_FILE} \
  get secrets -A
Error from server (Forbidden): secrets is forbidden: User "monitor" cannot list resource "secrets" in API group "" at the cluster scope
```

### setup.sh
```shell
#!/bin/sh

export KUBE_CLUSTER_ID=kubernetes
export KUBE_API_SERVER=https://192.168.77.77:6443
export KUBE_API_CA_KEY=/etc/kubernetes/pki/ca.key
export KUBE_API_CA_CRT=/etc/kubernetes/pki/ca.crt
export CLIENT_USR_NAME=monitor
export USER_NAME_SPACE=monitor

# =========== RBAC
export KUBE_RBAC_FILE=./rbac-${KUBE_CLUSTER_ID}-${CLIENT_USR_NAME}.yaml
echo "kind: Namespace                                " > ${KUBE_RBAC_FILE}
echo "apiVersion: v1                                 ">> ${KUBE_RBAC_FILE}
echo "metadata:                                      ">> ${KUBE_RBAC_FILE}
echo "  name: ${USER_NAME_SPACE}                     ">> ${KUBE_RBAC_FILE}
echo "  labels:                                      ">> ${KUBE_RBAC_FILE}
echo "    name: ${USER_NAME_SPACE}                   ">> ${KUBE_RBAC_FILE}
echo "---                                            ">> ${KUBE_RBAC_FILE}
echo "kind: ClusterRole                              ">> ${KUBE_RBAC_FILE}
echo "apiVersion: rbac.authorization.k8s.io/v1       ">> ${KUBE_RBAC_FILE}
echo "metadata:                                      ">> ${KUBE_RBAC_FILE}
echo "  namespace: ${USER_NAME_SPACE}                ">> ${KUBE_RBAC_FILE}
echo "  name: ${CLIENT_USR_NAME}                     ">> ${KUBE_RBAC_FILE}
echo "rules:                                         ">> ${KUBE_RBAC_FILE}
echo "- apiGroups:                                   ">> ${KUBE_RBAC_FILE}
echo "    - \"\"                                     ">> ${KUBE_RBAC_FILE}
echo "    - extensions                               ">> ${KUBE_RBAC_FILE}
echo "    - apps                                     ">> ${KUBE_RBAC_FILE}
echo "    - batch                                    ">> ${KUBE_RBAC_FILE}
echo "  resources:                                   ">> ${KUBE_RBAC_FILE}
echo "    - nodes                                    ">> ${KUBE_RBAC_FILE}
echo "    - services                                 ">> ${KUBE_RBAC_FILE}
echo "    - pods                                     ">> ${KUBE_RBAC_FILE}
echo "    - deployments                              ">> ${KUBE_RBAC_FILE}
echo "    - daemonsets                               ">> ${KUBE_RBAC_FILE}
echo "    - replicasets                              ">> ${KUBE_RBAC_FILE}
echo "    - replicationcontrollers                   ">> ${KUBE_RBAC_FILE}
echo "    - statefulsets                             ">> ${KUBE_RBAC_FILE}
echo "    - jobs                                     ">> ${KUBE_RBAC_FILE}
echo "    - cronjobs                                 ">> ${KUBE_RBAC_FILE}
echo "  verbs:                                       ">> ${KUBE_RBAC_FILE}
echo "    - get                                      ">> ${KUBE_RBAC_FILE}
echo "    - watch                                    ">> ${KUBE_RBAC_FILE}
echo "    - list                                     ">> ${KUBE_RBAC_FILE}
echo "---                                            ">> ${KUBE_RBAC_FILE}
echo "kind: ClusterRoleBinding                       ">> ${KUBE_RBAC_FILE}
echo "apiVersion: rbac.authorization.k8s.io/v1       ">> ${KUBE_RBAC_FILE}
echo "metadata:                                      ">> ${KUBE_RBAC_FILE}
echo "  namespace: ${USER_NAME_SPACE}                ">> ${KUBE_RBAC_FILE}
echo "  name: read-resources                         ">> ${KUBE_RBAC_FILE}
echo "subjects:                                      ">> ${KUBE_RBAC_FILE}
echo "- kind: User                                   ">> ${KUBE_RBAC_FILE}
echo "  name: ${CLIENT_USR_NAME}                     ">> ${KUBE_RBAC_FILE}
echo "  apiGroup: rbac.authorization.k8s.io          ">> ${KUBE_RBAC_FILE}
echo "roleRef:                                       ">> ${KUBE_RBAC_FILE}
echo "  kind: ClusterRole                            ">> ${KUBE_RBAC_FILE}
echo "  name: ${CLIENT_USR_NAME}                     ">> ${KUBE_RBAC_FILE}
echo "  apiGroup: rbac.authorization.k8s.io          ">> ${KUBE_RBAC_FILE}

kubectl apply -f ${KUBE_RBAC_FILE}

# =========== 认证
export KUBE_CONFIG_FILE=./kube-${KUBE_CLUSTER_ID}-${CLIENT_USR_NAME}.conf
openssl genrsa -out ./${CLIENT_USR_NAME}.key 2048
openssl req -new                                     \
 -subj "/CN=${CLIENT_USR_NAME}/O=${USER_NAME_SPACE}" \
 -key ./${CLIENT_USR_NAME}.key                       \
 -out ./${CLIENT_USR_NAME}.csr
openssl x509 -req -CAcreateserial -days 9999 \
  -CA ${KUBE_API_CA_CRT}                     \
  -CAkey ${KUBE_API_CA_KEY}                  \
 -in ./${CLIENT_USR_NAME}.csr -out ./${CLIENT_USR_NAME}.crt

# ============ 客户端配置
# 集群
# https://www.kubernetes.org.cn/doc-52
## --certificate-authority=/etc/kubernetes/pki/ca.crt # 设置kuebconfig配置文件中集群选项中的certificate-authority路径
## --embed-certs=true                                 # 设置kuebconfig配置文件中集群选项中的embed-certs开关
## --insecure-skip-tls-verify=false                   # 设置kuebconfig配置文件中集群选项中的insecure-skip-tls-verify开关
## --server="https://1.2.3.4:6443"                    # 设置kuebconfig配置文件中集群选项中的server
kubectl --kubeconfig=${KUBE_CONFIG_FILE} \
  config set-cluster ${KUBE_CLUSTER_ID}  \
    --server=${KUBE_API_SERVER}          \
    --embed-certs=true                   \
    --certificate-authority=${KUBE_API_CA_CRT} # --insecure-skip-tls-verify=true
    
# 用户
# https://www.kubernetes.org.cn/doc-54
## --client-certificate="": 设置kuebconfig配置文件中用户选项中的证书文件路径。
## --client-key="": 设置kuebconfig配置文件中用户选项中的证书密钥路径。
## --embed-certs=false: 设置kuebconfig配置文件中用户选项中的embed-certs开关。
## --password="": 设置kuebconfig配置文件中用户选项中的密码。
## --token="": 设置kuebconfig配置文件中用户选项中的令牌。
## --username="": 设置kuebconfig配置文件中用户选项中的用户名。
kubectl --kubeconfig=${KUBE_CONFIG_FILE}          \
  config set-credentials ${CLIENT_USR_NAME}       \
    --client-key=./${CLIENT_USR_NAME}.key         \
    --client-certificate=./${CLIENT_USR_NAME}.crt \
    --embed-certs=true                    # --username=admin --password=uXFGweU9l35qcif # --token=bearer_token

# 上下文
# https://www.kubernetes.org.cn/doc-53
## --cluster=""   # 设置kuebconfig配置文件中环境选项中的集群
## --user=""      # 设置kuebconfig配置文件中环境选项中的用户
## --namespace="" # 设置kuebconfig配置文件中环境选项中的命名空间
kubectl --kubeconfig=${KUBE_CONFIG_FILE}                   \
  config set-context ${CLIENT_USR_NAME}@${KUBE_CLUSTER_ID} \
    --cluster=${KUBE_CLUSTER_ID}                           \
    --user=${CLIENT_USR_NAME}                              \
    --namespace=${USER_NAME_SPACE}
    
kubectl --kubeconfig=${KUBE_CONFIG_FILE}                   \
  config set current-context ${CLIENT_USR_NAME}@${KUBE_CLUSTER_ID}

exit 0
```
