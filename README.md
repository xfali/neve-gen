# neve-gen
neve-gen是一个根据用户配置或解析指定数据库自动生成RESTful服务端及客户端代码的框架，支持：
* 生成CRUD Service interface
* 生成CRUD Service interface实现，包括空实现/实际数据库交互实现
* 生成数据库操作代理方法（gobatis mapper及proxy）
* 生成提供RESTful API服务，对应CRUD Service
* 生成RestClient（远程调用RESTful API服务），为CRUD Service interface的客户端实现
* 生成swagger文档

待完成
* 生成proto及grpc服务端

## 1 安装

使用命令安装：

```
go get github.com/xfali/neve-gen/cmd/neve-gen
```

## 2 使用

### 2.1、使用内置模板
```
neve-gen -c config.yaml -f value.yaml -o TARGET_DIR
```

### 2.2、指定模板目录

```
neve-gen -c config.yaml -f value.yaml -t TEMPLATE_DIR -o TARGET_DIR
```

## 3 说明
### 3.1 配置（config.yaml）
config.yaml指定代码生成的规则及配置，包括
1. 是否生成swagger
2. 是否生成gobatis
3. 是否生成RestClient
4. 指定扫描的数据库及相关信息
5. 其他请参考[示例](example/conf.yaml)

请根据项目需要修改相关值的内容。

### 3.2 值定义（value.yaml）
value.yaml指定生成代码所需的值，neve-gen会根据配置的值自动解析和生成相应的代码，包括
1. 作者信息 
2. 应用名称、版本号
3. web服务端口
4. 数据库配置
5. 返回结果的结构
6. 自定义的数据结构
7. 其他请参考[示例](example/value.yaml)

请根据项目需要修改相关值的内容。

### 3.3 模板（template）
template是代码生成的核心，包含代码生成规则的描述及使用的模板：
1. 模板的描述文件[layout.yaml](templates/layout.yaml)，该文件描述了模板的组织形式，约定了neve-gen如何使用模板和生成代码
2. 各类以tmpl为后缀的模板文件，neve-gen将结合config.yaml，value.yaml来自动生成目的代码文件。

模板须经过严格的测试和验证，请不要随意修改模板相关内容。

（neve-gen内置模板生成的目的工程为基于neve框架的web服务应用。）

### 3.4 自定义模板开发
待续...