[bearychat_ifttt](https://github.com/flymzero/bearychat_ifttt) 使用帮助

## 更新

> 去除email数据绑定，一个对象只绑定**名称**和**key**，名称唯一

## 原理说明
> 利用ifttt的webhooks服务通过**触发词**来执行你指定的操作（3个参数可提供）
> 
> ifttt相关文章: [链接](https://sspai.com/post/39243?utm_source=weibo&utm_medium=sspai&utm_campaign=weibo&utm_content=ifttt&utm_term=jiaocheng)
> 
> ifttt获取key: [链接](http://maker.ifttt.com/)

## 对象数据绑定操作

把自己或者别人（统一称为对象）先绑定对应的ifttt的key执行操作
一个对象包含**名称**和**key**两个字段 （名称唯一）

> **-h**    获取使用帮助

> **-ls**   列出所有自己及自己设置对象的信息
> 
 ![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/721541038190_.pic.jpg)

> **-s [-m]  n:名称  k:ifttt上的key  **   设置自己或者对象的信息(n必填唯一，-m则表示设置自己信息)
- 例1 : -s -m n:我 k:xxxxxx  设置自己的名称和ifttt的key
- 例2 : -s n:老婆 k:oooooo 设置一个对象：老婆及触发的key
  
![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/681541038185_.pic.jpg)
> **-d 名称**  删除你对象中对应名称的这个对象
- 例 : -d 老婆 把老婆这个对象删了
  
  ![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/671541038184_.pic.jpg)


## 触发操作

> **$触发词  [n:名称]  [v1:xx]  [v2:xx]  [v3:xx] ** 
> 
对这个名称对象(n不填就是自己，v1,v2,v3都是可选的)，进行触发操作,并传输可选的3个参数，~~当存在引用附件时以附件的url优先作为v3的值~~,贝洽的文件在未登录的情况下无法访问链接,所以无法转存!!!

需在"对象"手机ifttt上创建对应的Applet

**基本上ifttt上能创造的东西(通知,发邮件,文件转存,发微博....),你都可以让这个机器人代劳, 话说ifttt不开放共享applet也真是坑**


## ifttt配置
- 开通webhooks服务
  
  ![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/841541038268_.pic.jpg)

- 获取ifttt的webhook的key
  
  ![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/881541038275_.pic.jpg)


## 示例

> **发邮件给自己**


> 设置触发词

![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/861541038273_.pic.jpg)
![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/761541038224_.pic.jpg)

>选择触发服务email

![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/791541038229_.pic.jpg)

> 设置参数

![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/751541038223_.pic.jpg)

> 保存

**通过bearychat触发操作**

![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/741541038220_.pic.jpg)

> 查看

![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/771541038227_.pic.jpg)


**远程通知**

> 命令: $bc_notice n:老婆 v1:老婆大人查收 v2:今天又要加班,你先睡吧 v3:https://dwz.cn/UEPUo3bC
![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/%20notice.jpg)

> 发微博

![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/731541038190_.pic.jpg)
![示例](https://raw.githubusercontent.com/flymzero/bearychat_ifttt/master/imgs/701541038186_.pic.jpg)

## 还有更多等你自己拓展