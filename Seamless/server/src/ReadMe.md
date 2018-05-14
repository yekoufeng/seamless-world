# 目录说明

目录        | 说明
------------|----------------------------------
Center      | Center服务器
common      | 公共包, 一些protobuf协议定义、其他功能定义等
DataCenter  | DataCenter服务器
datadef     | datadef包, 保存到mysql数据库的一些结构体定义
db          | db包, 访问redis的一些操作类
entitydef   | entitydef包, 实体定义，由 res/entitydef 生成
excel       | excel包, excel配置表读取，由 res/excel 生成
Gateway     | Gateway服务器
generator   | entitydef 生成工具，res/entitydef/*.json 生成 src/datadef/*.go
idip        | idip包，Center服务器跟IDIP服务器交互的逻辑
IDIPServer  | IDIP服务器
ImportCdkey | 导入CDKEY的工具
Lobby       | Lobby服务器
Login       | Login服务器
Match       | Match服务器
msdk        | msdk包
Pay         | Pay服务器
protoMsg    | protoMsg包，由 game.proto 生成的 protobuf go文件
Room        | Room服务器
stress      | 压力测试工具
utility     | 导入黑名单工具
vendor      | 第3方库
zeus        | zeus包，服务器底层框架库, 详见 [zeus/README.md](zeus/README.md)

