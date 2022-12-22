
const userDeviceList = await getUserDeviceList();

const userDOMList = document.createElement("ul");

for(const user of Object.keys(userDeviceList)){
	console.log(user);
	const DOMUser = document.createElement("li");
	userDOMList.append(DOMUser);
	const DOMUserName = document.createElement("p");
	DOMUserName.innerText = user;
	DOMUser.append(DOMUserName);
	const deviceDOMList = document.createElement("ul");
	DOMUser.append(deviceDOMList);
	for(const device of userDeviceList[user]){
		const DOMDevice = document.createElement("li");
		deviceDOMList.append(DOMDevice);
		
		const deviceCheckbox = document.createElement("input");
		deviceCheckbox.type = "checkbox";
		DOMDevice.append(deviceCheckbox);

		const DOMDeviceName = document.createElement("p");
		DOMDeviceName.innerText = device;
		DOMDevice.append(DOMDeviceName);
	}
}

document.body.append(userDOMList);

async function getUserDeviceList() {
	const userDeviceList = {};

	const users = (await fetchJson("/api/0/list")).results;
	users.sort();
	for(const user of users){
		const devices = (await fetchJson(`/api/0/list?user=${encodeURIComponent(user)}`)).results;
		devices.sort();
		userDeviceList[user] = devices;
	}

	return userDeviceList;
}

async function fetchJson(url){
	return (await (await fetch(url)).json());
}

