pragma solidity ^0.4.17;

contract Donate {
  address public owner;//合约拥有者地址

  //受捐赠人结构体
  struct Donee {
        uint id;//序号
        string name;//姓名
        string age; //年龄
        string city;//所在城市
        string reason;//受捐赠的原因
        int receiveAmount;//已接受捐款数目
  }
  
  //捐赠记录
  struct Record {
        address donor; //
        string  doneeName;
        int amount;
  }
  
 
  //捐赠人捐款总金额 (捐赠人地址->已捐款金额)
  mapping(address => int) public donors;

  //捐款记录
  Record[] public recordList;
  
  //受捐赠人map (序号->受捐赠人)
  mapping (uint => Donee) public doneeMapping;

  //记录受捐赠者id
  uint public doneeCount;

  constructor () public {
    owner = msg.sender;
    initDonees();
  }

  function initDonees ()  private {
    addDonee("张三", "60", "西安", "心脏病");
    addDonee("李四", "47", "拉萨", "胃癌");
    addDonee("王五", "12", "南京", "白血病");
  }


  function addRecord (address _donorAdrress, string _doneeName, int _amount) private {
    recordList.push(Record(_donorAdrress, _doneeName, _amount));
  }
  

  //添加受捐赠人
  function addDonee (string _doneeName, string _doneeAge, string _doneeCity, string _doneeReason) public {
    doneeCount++;
    doneeMapping[doneeCount] = Donee(doneeCount, _doneeName, _doneeAge, _doneeCity, _doneeReason, 0);
  }
  

  //捐款
  function donateTo (uint _doneeID, int _amount) public payable {
    //starID要符合
    require(_doneeID > 0 && _doneeID <= doneeCount);
    //捐款数目要大于0
    require (_amount > 0);

    //受捐赠人map (序号->受捐赠人)
    doneeMapping[_doneeID].receiveAmount = doneeMapping[_doneeID].receiveAmount + _amount;
    
    //累加
    donors[msg.sender] = donors[msg.sender] + _amount;

    //添加记录
    addRecord(msg.sender, doneeMapping[_doneeID].name, _amount);

  }
}
