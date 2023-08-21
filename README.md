## kex ##

kubectl extend
- show: show pod list
- enter: use pods select list to enter pod
- change: set config to ~/.kube/config
- describe: select pod to describe
- searchlog: TODO searh multi pod log

## Install ##
1. install kubectl 
2. go build .
3. mv kex /usr/local/bin/

## Usage ##
- single cluster config
mkdir ~/.kube
mv you config file to ~/.kube/config

```
$kex show
```

- multi cluster config
1. mkdir -p ~/.kube/configs
2. add you cluster config to ~/.kube/configs/

```
$kex change
$kex show
```

## completion ##
- todo
