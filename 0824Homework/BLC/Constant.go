package BLC

//常量值
const DBNameYS = "blockchain_%s.db" //数据库的名字
const BlockBucketNameYS = "blocks" //定义bucket
const targetBitYS = 8              // 挖矿难度(256位Hash里面前面至少要有16个零)
const UTXOSetBucketNameYS = "utxoset"
const walletsFileYS = "Wallets_%s.dat"//存储钱包数据的本地文件名
const txPollFileYS = "TxsPool_%s.dat"//本地交易池

//网络相关
const NODE_VERSIONYS = 1    //版本
const COMMAND_LENGTHYS = 12 //命令长度[]byte
const BLOCK_TYPEYS = "BLOCK_TYPE"
const TX_TYPEYS = "TX_TYPE"

//具体的命令
const COMMAND_VERSIONYS = "version"
const COMMAND_GETBLOCKSYS = "getblocks"
const COMMAND_INVYS = "inv"
const COMMAND_GETDATAYS = "getdata"
const COMMAND_BLOCKDATAYS = "blockdata"
const COMMAND_TXSYS = "transactions"
const COMMAND_REQUIREMINEYS = "requiremine"
const COMMAND_VERIFYBLOCKYS  = "verifyblock"

//钱包
const versionYS = byte(0x00)
const addressCheckSumLenYS = 4


