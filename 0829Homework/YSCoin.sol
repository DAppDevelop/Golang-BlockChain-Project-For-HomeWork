pragma solidity ^0.4.24;

contract EIP20{
    function totalSupply() view returns (uint256 totalSupply);
    function balanceOf(address _owner) view returns (uint256 balance);
    function transfer(address _to, uint256 _value) returns (bool success);
    function approve(address _spender, uint256 _value) returns (bool success);
    function transferFrom(address _from, address _to, uint256 _value) returns (bool success);
    function allowance(address _owner, address _spender) view returns (uint256 remaining);
    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}

contract YSCoin is EIP20{

    string public name;//ERC20代币的名字
    string public symbol;//ERC20代币的符号
    uint8 public decimals;//支持几位小数点后几位

    mapping (address=>uint256) balances;//记录每个地址拥有的代币数目
    mapping(address=>mapping(address=>uint256) )allowances;//记录每个地址允许某地址从自己账户转走代币数目

    function YSCoin(){
        name = "YSCion";
        symbol = "YSC";
        decimals = 2;
        balances[msg.sender] = totalSupply();
    }

    //发币总量
    function totalSupply() view returns (uint256 totalSupply) {
        return 10000000;
    }

    //查询某个地址拥有的代币数量
    function balanceOf(address _owner) view returns (uint256 balance){
        return balances[_owner];
    }

    //调用transfer函数将自己的token转账给_to地址，_value为转账个数
    function transfer(address _to, uint256 _value) returns (bool success){

        require(_value > 0 && balances[msg.sender] >= _value && balances[_to] + _value > balances[_to]);
        //修改发送方和接收方的代币数量
        balances[_to] += _value;
        balances[msg.sender] -= _value;

        //触发事件
        Transfer(msg.sender, _to, _value);
        return true;
    }

    //批准_spender账户从自己的账户转移_value个token。可以分多次转移
    function approve(address _spender, uint256 _value) returns (bool success){

        require (balances[msg.sender] > _value);

        allowances[msg.sender][_spender] = _value;
        //触发事件
        Approval(msg.sender,_spender,_value);

        return true;
    }

    //与approve搭配使用，approve批准之后，调用transferFrom函数来转移token。  B可以从A转移代币到地址C
    function transferFrom(address _from, address _to, uint256 _value) returns (bool success){
        //获取可以转移的数量
        uint256 allowan = allowances[_from][msg.sender];

        require(_value > 0 && balances[_from] >= _value && allowan >= _value  && balances[_to]+_value > balances[_to]);

        allowances[_from][msg.sender] -= _value;

        balances[_from] -= _value;
        balances[_to] += _value;

        Transfer(_from,_to,_value);

        return true;

    }


    //返回_spender还能从_owner提取token的个数
    function allowance(address _owner, address _spender) view returns (uint256 remaining){
        return allowances[_owner][_spender];
    }

}
