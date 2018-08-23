package BLC

const targetBitYS = 8              // 挖矿难度(256位Hash里面前面至少要有targetBit个零)
const DBNameYS = "blockchainYS.db" // 数据库名字
const BlockBucketNameYS = "blocks" // 表的名字
const UTXOSetBucketNameYS = "utxoset"
const walletsFileYS = "Wallets.dat" //存储钱包数据的本地文件名


//钱包
const versionYS = byte(0x00)
const addressCheckSumLenYS = 4
