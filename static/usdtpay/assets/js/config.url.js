//
global_requestAddress = "";

//登录接口
global_requestAddressJs_login = global_requestAddress+ "/v2/login";

//获取菜单
global_requestAddressJs_getMenus = global_requestAddress+ "/v2/getMenus";

//获取订单/回调/修改备注
global_requestAddressJs_topUpOrder = global_requestAddress+ "/order/topUpOrder";

//获取地址管理/更新所有余额/单个余额/资金归集
global_requestAddressJs_toAddress = global_requestAddress+ "/address/toAddress";







// 国际时间-时区设置
var currTimeZoneValue = 5










const timeArr = [{
	name:'Asia/Shanghai',
	value:+8,
},{
	name:'Asia/Kolkata',
	value:+5.5,
},{
	name:'America/Mexico_City',
	value:-6,
},{
	name:'America/Bogota',
	value:-5,
},{
	name:'Pakistan',
	value:5,
}]

var rechargeWithdrawTypeList = [{
	name:'USDT',
	value:1
},{
	name:'BPAY',
	value:2
},{
	name:'lrPay',
	value:3
},{
	name:'wowPay',
	value:4
},]

var anotherPayTypeList = [{
	name:'本地代付',
	value:1
},{
	name:'BPAY',
	value:2
},{
	name:'lrPay',
	value:3
},{
	name:'wowPay',
	value:4
},]


var param={};

let currTimeZoneNameStr
//
// //获取配置数据
// $.ajax({
// 	url: global_requestAddress + global_requestAddress_js_systemParameter+"?action=select",
// 	headers:{
// 		token:$.cookie('tokenMybUP')
// 	},
// 	dataType: 'json',
// 	type: 'post',
// 	data: param,
// 	success: function (dataArray) {
// 		if (dataArray.code != 2000) {
// 			return false;
// 		}
// 		var returnDataArray = dataArray.result
// 		let currTimeZoneName = returnDataArray.TimeZone
// 		currTimeZoneNameStr = returnDataArray.TimeZone
//
// 		if(returnDataArray.TimeZone){
// 			// console.log("currTimeZoneValue111",currTimeZoneNameStr,returnDataArray)
// 			timeArr.forEach((item)=>{
// 				if(item.name === currTimeZoneName){
// 					currTimeZoneValue = item.value
// 				}
// 			})
// 		}
//
// 	}
// })




//-------------


var getRootPath_webStr = getRootPath_web();

//获取目录路径方法
function getRootPath_web() {

		//获取当前网址，如： http://localhost:8888/eeeeeee/aaaa/vvvv.html
		var curWwwPath = window.document.location.href;
		//获取主机地址之后的目录，如： uimcardprj/share/meun.jsp
		var pathName = window.document.location.pathname;
		var pos = curWwwPath.indexOf(pathName);
		//获取主机地址，如： http://localhost:8888
		var localhostPaht = curWwwPath.substring(0, pos);
		//获取带"/"的项目名，如：/abcd
		var projectName = pathName.substring(0, pathName.substr(1).indexOf('/') + 1);

		// return (localhostPaht + projectName);


		// console.log("当前网址:"+curWwwPath);
		// console.log("主机地址后的目录:"+pos+"----"+pathName);
		// console.log("主机地址:"+localhostPaht);
		// console.log("项目名:"+projectName);


		return projectName;
}



//时间戳转日期时间型工具类
function formatDateTime(inputTime) {
	var date = new Date(inputTime);
	var y = date.getFullYear();
	var m = date.getMonth() + 1;
	m = m < 10 ? ('0' + m) : m;
	var d = date.getDate();
	d = d < 10 ? ('0' + d) : d;
	var h = date.getHours();
	h = h < 10 ? ('0' + h) : h;
	var minute = date.getMinutes();
	var second = date.getSeconds();
	minute = minute < 10 ? ('0' + minute) : minute;
	second = second < 10 ? ('0' + second) : second;
	return y + '-' + m + '-' + d+' '+h+':'+minute+':'+second;
}


function toDecimal2(x) {//金额处理两位小数点
	var f = parseFloat(x);
	if (isNaN(f)) {
		return false;
	}
	var f = Math.round(x*100)/100;
	var s = f.toString();
	var rs = s.indexOf('.');
	if (rs < 0) {
		rs = s.length;
		s += '.';
	}
	while (s.length <= rs + 2) {
		s += '0';
	}
	return s;
}


/**
 * 数字转整数 如 100000 转为10万
 * @param {需要转化的数} num
 * @param {需要保留的小数位数} point
 */
function tranNumber(num, point) {



	let numStr = num.toString()

	// console.log(numStr.length);
	// 一万以内直接返回
	if (numStr.length <=4) {
		return numStr;
	}
	//大于6位数是十万 (以10W分割 10W以下全部显示)
	else if (numStr.length > 4) {
		let decimal = numStr.substring(numStr.length - 4, numStr.length - 4 + point)
		// return parseFloat(parseInt(num / 10000) + ‘.’ + decimal) + ‘万’;
		return parseFloat(parseInt(num / 10000) + '.' + decimal) + '万';
	}
}




//验证是否为数字
function isNumber(value) { //验证是否为数字

	var patrn = /^(-)?\d+(\.\d+)?$/;

	if (patrn.exec(value) == null || value == "") {
		return false

	} else {
		return true

	}

}



// 北京是getZoneTime(8),纽约是getZoneTime(-5),班加罗尔是getZoneTime(5.5). 偏移值是本时区相对于格林尼治所在时区的时区差值
function getZoneTime(offset){
	// 取本地时间
	var localtime = new Date();
	// 取本地毫秒数
	var localmesc = localtime.getTime();
	// 取本地时区与格林尼治所在时区的偏差毫秒数
	var localOffset = localtime.getTimezoneOffset() * 60000;
	// 反推得到格林尼治时间
	var utc = localOffset + localmesc;
	// 得到指定时区时间
	var calctime = utc + (3600000*offset);
	var nd = new Date(calctime);
	return nd.toDateString()+" "+nd.getHours()+":"+nd.getMinutes()+":"+nd.getSeconds();
}


// 我们从后台拿到时区与时间戳，要转换为对应的时区时间。
// 可以在全局过滤器中写一个方法：
function getLocalTime(i,time){
	// i为传入的时区，东八区传8，东七传7
	// time为传入的时间戳,如1619712000000，这两都是从后台拿到的数据
	// 如果需要当前的时间戳（1970年一月一日到现在的秒数）
	// let date = new Date().getTime()

	// 得到本地时间与GMT时间的时间偏移差
	let offset = new Date().getTimezoneOffset() * 60 * 1000
	// 后台给的时间戳与offset相加得到现在的格林尼治时间
	let utcTime = time + offset
	// 得到正确的时间 如果要转成时间戳就加getTime()方法
	return new Date(utcTime + 60 * 60 * 1000 * i)
}

/**
 * 获得当前日期（年-月-日）
 */
function getCurrDate() {
	var date = new Date();
	var sep = "-";
	var year = date.getFullYear(); //获取完整的年份(4位)
	var month = date.getMonth() + 1; //获取当前月份(0-11,0代表1月)
	var day = date.getDate(); //获取当前日
	if (month <= 9) {
		month = "0" + month;
	}
	if (day <= 9) {
		day = "0" + day;
	}
	var currentdate = year + sep + month + sep + day;
	return currentdate;
}



function get_unix_time(dateStr) {
	var newstr = dateStr.replace(/-/g, '/');
	var date = new Date(newstr);
	var time_str = date.getTime().toString();
	return time_str.substr(0, 10);
}
