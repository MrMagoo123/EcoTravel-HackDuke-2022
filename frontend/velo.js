// Velo API Reference: https://www.wix.com/velo/reference/api-overview/introduction

const axios = require('axios');
import wixAnimations from 'wix-animations';

$w.onReady(() => {

	let urls = {
		"kiwi.com": "kiwi.com",
		"amtrak": "https://www.amtrak.com/acela-train",
		"megabus": "https://us.megabus.com/",
		"flix": "https://www.flixbus.com/bus/united-states",
		"go bus": "https://www.gobuses.com/",
	};

	const numResultsText = $w('#text28')
	numResultsText.hide();

	let doneLoading = false;

	function startLoading() {
		let rollIn = {
		"duration":   500,
		"delay": 0,
		"direction": "right"
		};

		let rollOut = {
		"duration":   500,
		"delay": 0,
		"direction": "left"
		};

		if (doneLoading) {
			return;
		}

		const searchButton = $w("#button4");
		searchButton.hide("roll", rollIn)  .then( ( ) => {
			searchButton.show("roll", rollOut).then( () => startLoading())
		} );
	}


	$w("#addressInput1").required = true;
	$w("#addressInput2").required = true;
	$w("#addressInput1").placeholder = "Enter starting address";
	$w("#addressInput2").placeholder = "Enter destination address";
	$w("#datePicker2").required = true;

	var priceLow = $w('#rangeSlider2').value[0];
	var priceHigh = $w('#rangeSlider2').value[1];
	var maxTravelTime = $w('#timePicker2').value;
	var sortCriteria = $w('#dropdown2').value;

	$w('#button3').onClick(() => {
		priceLow = $w('#rangeSlider2').value[0];
		priceHigh = $w('#rangeSlider2').value[1];
		maxTravelTime = $w('TimePicker').value;
		maxTravelTime = maxTravelTime.substring(0, maxTravelTime.length - 7);
		console.log("price low changed to: ", priceLow);
		console.log("price high changed to: ", priceHigh);
		console.log("max travel time changed to: ", maxTravelTime);
		var newData = [];
		for (var i = 0; i < allData.length; i++) {
			console.log("duration: ", allData[i].overallData.duration);
			if (allData[i].overallData.duration <= maxTravelTime && allData[i].overallData.price.substring(1) >= priceLow && allData[i].overallData.price.substring(1) <= priceHigh) {
				newData.push(allData[i]);
			}
		}
		$w("#repeater1").data = newData;
	});

	var allData = [];

	$w.onReady(() => {
		$w("#repeater1").onItemReady( ($item, itemData) => {
			console.log("overall data: ", itemData.overallData.price);
			$item("#table1").rows = itemData.data;
			$item("#text5").text = itemData.overallData.departure;
			$item("#text8").text = itemData.overallData.arrival;
			$item("#text13").text = itemData.overallData.price;
			$item("#text12").text = itemData.overallData.ecopoints;
			$item("#text26").text = itemData.overallData.duration;
			$item("#button5").link = itemData.overallData.bookingURL;
		});
	});
	
	$w("#repeater1").data = allData;


	$w('#dropdown2').onChange(() => {
		sortCriteria = $w('#dropdown2').value;
		console.log("sort criteria changed to: ", sortCriteria);
		if (sortCriteria == "Price") {
			console.log("sort by eco");
			allData = allData.sort(compareByPrice);
		} else if (sortCriteria == "Time") {
			allData = allData.sort(compareByTime);
		} else if (sortCriteria == "Eco") {
			console.log("sort by eco");
			allData = allData.sort(compareByEco);
		} else {
			console.log("set to original list")
			allData.sort(compareByEco);
		}
		console.log("alldata: ", allData);
		$w("#repeater1").data = allData;
	});

	function compareByTime(a, b) {
		// TODO
		return a.overallData.durationNum - b.overallData.durationNum;
	}

	function compareByPrice(a, b) {
		return a.overallData.price.substring(1) - b.overallData.price.substring(1);
	}

	function compareByEco(a, b) {
		return b.overallData.ecopoints - a.overallData.ecopoints;
	}

	$w('#button4').onClick(() => {
		doneLoading = false
		startLoading()
		console.log("data: ", $w("#repeater1").data);
		console.log("start: ", $w("#addressInput1").value);
		console.log("end: ", $w("#addressInput2").value);
		console.log("date: ", $w('#datePicker2').value);
		console.log("price: ", $w('#rangeSlider2').value);

		const start = $w("#addressInput1").value;
		const end = $w("#addressInput2").value;
		const date = $w("#addressInput1").value;

		const requestBody = {
			"startLocation": start,
			"endLocation": end,
			"date": date,
		};

		const payload = JSON.stringify(requestBody);

		const headers = {
			"Content-Type": "application/json"
		};

		var receivedData = [];

		axios.post("https://api.ecotravel.tech/api/createQuery", payload, { headers })
		.then((response) => {
			console.log("response: ", response);
			if (!response.data) {
				var slug = response.slug;
			} else {
				receivedData = response.data;
			}
			console.log("new list data: ", receivedData);
			var newData = [];
			for (var i = 1; i <= receivedData.length; i++) {
                var curRoute = receivedData[i - 1];
                var duration = curRoute.duration;
                var hours = Math.trunc(duration / 3600);
                duration -= hours * 3600;
                var mins = Math.trunc(duration / 60);
                var durationString = hours + ":";
                if (hours < 10) {
                    durationString = "0" + durationString;
                }
                if (mins < 10) {
                    durationString = durationString + "0" + mins;
                } else {
                    durationString = durationString + mins;
                }
                var ecopoints = 10;
                if (curRoute.description.includes("fly") || curRoute.description.includes("Fly")) {
                    ecopoints += 5;
                } 
                if (curRoute.description.includes("Drive") || curRoute.description.includes("drive")) {
                    ecopoints += 15;
                }
                if (curRoute.description.includes("bus") || curRoute.description.includes("Bus") || curRoute.carrier.includes("greyhound") || curRoute.carrier.includes("Greyhound") || curRoute.description.includes("Peter Pan Bus Lines")) {
                    ecopoints += Math.trunc((234*(hours + (mins/60))*13)*0.01);
                }
                
                if (curRoute.carrier.includes("amtrak") || curRoute.carrier.includes("Amtrak") || curRoute.description.includes("Train") || curRoute.description.includes("train") || curRoute.type.includes("bus") || curRoute.type.includes("Bus")) {
                    ecopoints += Math.trunc((112*(hours + (mins/60))*41)*0.01);
                }
                if (ecopoints < 0) {
                    ecopoints = 0
                }

				var startTime = new Date(1000 * curRoute.startTime);
				var endTime = new Date(1000 * curRoute.endTime);
				var startTimeString = "";
				var startHours = startTime.getHours();
				if (startHours < 10) {
					startTimeString += "0" + startHours + ":";
				} else {
					startTimeString += startHours + ":";
				}
				var startMins = startTime.getMinutes();
				if (startMins < 10) {
					startTimeString += "0" + startMins;
				} else {
					startTimeString += startMins;
				}

				var endTimeString = "";
				var endHours = endTime.getHours();
				if (endHours < 10) {
					endTimeString += "0" + endHours + ":";
				} else {
					endTimeString += endHours + ":";
				}
				var endMins = endTime.getMinutes();
				if (endMins < 10) {
					endTimeString += "0" + endMins;
				} else {
					endTimeString += endMins;
				}

				console.log("start: ", startTimeString);
				console.log("end: ", endTimeString);

				const lowercaseCarrier = curRoute.carrier.toLowerCase();

				let bookingURL;
				for (let site of Object.keys(urls)) {
					if (lowercaseCarrier.includes(site)) {
						bookingURL = urls[site];
					}
				}

				var newTrip = {
					"_id": i.toString(),
					"overallData": {
						"bookingURL": bookingURL,
						"departure": startTimeString,
						"arrival": endTimeString,
						"price": "$" + curRoute.price.toString(),
						"ecopoints": ecopoints.toString(),
						"duration": durationString,
						"description": curRoute.description,
						"durationNum": curRoute.duration,
					},
					"data": [
						{
						"company": curRoute.carrier,
						"price": "$" + curRoute.price.toString(),
						"starting_location": curRoute.startLocation,
						"departure": startTimeString,
						"destination": curRoute.endLocation,
						"arrival": endTimeString
						}
					]
				}
				newData.push(newTrip);
			}
			allData = newData;
			allData.sort(compareByEco);
			console.log("final new data: ", allData);

			$w("#repeater1").data = newData;

			doneLoading = true;
			numResultsText.text = `${$w("#repeater1").data.length.toString()} results`;
			numResultsText.show();

			setTimeout(() => {
      			$w("#anchor1").scrollTo()
      		}, 100);
		})
		.catch((error) => {
			console.log(error)
		});
	});
});
