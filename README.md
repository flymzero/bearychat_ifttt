[bearychat_ifttt](https://github.com/flymzero/bearychat_ifttt) 使用帮助
**对象操作**
> **-h**    获取使用帮助

> **-ls**   列出所有自己及自己设置对象的信息

> **-s [-m]  n:昵称  k:ifttt上的key  [e:邮箱]**   设置自己或者对象的信息(n必填，-m则表示设置自己信息)
- 例1 : -s -m n:我 k:xxxxxx  设置自己的昵称和ifttt的key
- 例2 : -s n:老婆 k:oooooo e:xxx@gmail.com 设置一个对象：老婆及触发的key
> **-d 昵称**  删除你对象中对应昵称的这个对象
- 例 : -d 老婆 把老婆这个对象删了

**触发操作**
> **$触发词  [n:昵称]  [v1:xx]  [v2:xx]  [v3:xx] ** 
对这个昵称对象(n不填就是自己，v1,v2,v3都是可选的)，进行触发操作,并传输可选的3个参数，~~当存在引用附件时以附件的url优先作为v3的值~~,贝洽的文件在未登录的情况下无法访问链接,所以无法转存!!!

需在"对象"手机ifttt上创建对应的Applet

ifttt相关文章: [链接](https://sspai.com/post/39243?utm_source=weibo&utm_medium=sspai&utm_campaign=weibo&utm_content=ifttt&utm_term=jiaocheng)
ifttt获取key: [链接](http://maker.ifttt.com/)

**基本上ifttt上能创造的东西(通知,发邮件,文件转存,发微博....),你都可以让这个机器人代劳, 话说ifttt不开放共享applet也真是坑**

**示例**
> 远程通知

命令: $bc_notice n:老婆 v1:老婆大人查收 v2:今天又要加班,你先睡吧 v3:https://dwz.cn/UEPUo3bC
![示例](https://s1.ax1x.com/2018/11/01/iWUU2D.jpg)