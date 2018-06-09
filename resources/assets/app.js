var webSocket;
var rpmGauge;
var gearElement = document.getElementById("gear-text");
var kphElement = document.getElementById("kph-text");
var progressElement = document.getElementById("progress");
var distanceElement = document.getElementById("distance");
var rpmGaugeCanvas = document.getElementById('rpm');

var rpmGaugeOptions = {
    angle: -0.3, // The span of the rpmGauge arc
    lineWidth: 0.2, // The line thickness
    pointer: {
        length: 0, // The radius of the inner circle
        strokeWidth: 0, // The rotation offset
        color: '#522A08' // Fill color
    },
    radiusScale: 1, // Relative radius
    limitMax: false,     // If false, max value increases automatically if value > maxValue
    limitMin: false,     // If true, the min value of the rpmGauge will be fixed
    strokeColor: 'rgba(238, 238, 238, 0.4)',  // to see which ones work best for you
    generateGradient: true,
    highDpiSupport: true,     // High resolution support
    percentColors: [[0.0, "#a9d70b"], [0.60, "#a9d70b"], [0.75, "#f9c802"], [1.0, "#ff0000"]],
    renderTicks: {
        divisions: 10,
        divWidth: 0.5,
        divLength: 0.2,
        divColor: "#333333"
    },
    staticLabels: {
        font: "10px sans-serif",  // Specifies font
        labels: [0, 1000, 2000, 3000],  // Print labels at these values
        color: "#000000",  // Optional: Label text color
        fractionDigits: 0  // Optional: Numerical precision. 0=round off.
    },
};
rpmGauge = new Gauge(rpmGaugeCanvas).setOptions(rpmGaugeOptions); // create sexy rpmGauge!
rpmGauge.maxValue = 3000; // set max rpmGauge value
rpmGauge.setMinValue(0);  // Prefer setter over rpmGauge.minValue = 0
rpmGauge.animationSpeed = 2; // set animation speed (32 is default value)
rpmGauge.set(0); // set actual value


var pedalPositions = new SmoothieChart({
    minValue: 0,
    maxValue: 1,
    labels: {disabled: true},
    grid:{fillStyle:'transparent',millisPerLine:10000,verticalSections:0}
});
pedalPositions.streamTo(document.getElementById("pedalcanvas"), 0);
var throttleLine = new TimeSeries();
var brakeLine = new TimeSeries();
var clutchLine = new TimeSeries();
pedalPositions.addTimeSeries(throttleLine, {
    strokeStyle: 'rgb(0, 255, 0)',
    fillStyle: 'rgba(0, 255, 0, 0.3)',
    lineWidth: 3
});
pedalPositions.addTimeSeries(brakeLine, {
    strokeStyle: 'rgb(255, 0, 0)',
    fillStyle: 'rgba(255, 0, 0, 0.3)',
    lineWidth: 3
});
pedalPositions.addTimeSeries(clutchLine, {
    strokeStyle: 'rgb(0, 0, 255)',
    fillStyle: 'rgba(0, 0, 255, 0.3)',
    lineWidth: 3
});


function renderGauge(telemetryData) {
    rpmGauge.maxValue = (Math.round((telemetryData.maxRpm * 10) / 1000) * 1000) + 1000;
    rpmGauge.set(telemetryData.rpm * 10);

    var labels = [];
    for (var i = 0; i <= rpmGauge.maxValue; i = i + 1000) {
        labels.push(i);
    }

    rpmGauge.setOptions({
        staticLabels: {
            font: "12px 'Squada One', cursive",  // Specifies font
            labels: labels,  // Print labels at these values
            color: "#ACACAC",  // Optional: Label text color
            fractionDigits: 0  // Optional: Numerical precision. 0=round off.
        },
        renderTicks: {
            divisions: rpmGauge.maxValue / 1000,
            divWidth: 0.5,
            divLength: 0.2,
            divColor: "#333333"
        },
    });
}

function renderTexts(telemetryData) {
    kphElement.innerHTML = Math.round(telemetryData.speed * 3.6);
    distanceElement.innerHTML = (telemetryData.lapDistance / 1000).toFixed(1) + " km";
    progressElement.innerHTML = Math.round(telemetryData.lapDistance / telemetryData.trackLength * 100) + " %";

    var gearText = "N";
    switch (telemetryData.gear) {
        case 1:
            gearText = "1";
            break;
        case 2:
            gearText = "2";
            break;
        case 3:
            gearText = "3";
            break;
        case 4:
            gearText = "4";
            break;
        case 5:
            gearText = "5";
            break;
        case 6:
            gearText = "6";
            break;
        case 7:
            gearText = "7";
            break;
        case 10:
            gearText = "R";
            break;
        default:
            gearText = "N";
    }

    gearElement.innerHTML = gearText;
}

function renderPedalPositions(telemetryData) {
    throttleLine.append(new Date().getTime(), telemetryData.throttlePosition);
    brakeLine.append(new Date().getTime(), telemetryData.brakePosition);
    clutchLine.append(new Date().getTime(), telemetryData.clutchPosition);
}

function render(telemetryData) {
    renderGauge(telemetryData);
    renderTexts(telemetryData);
    renderPedalPositions(telemetryData);
}

function connect(host) {
    webSocket = new WebSocket(host);
    webSocket.onopen = function (evt) {
        console.log("Websocket opened");
    };
    webSocket.onclose = function (evt) {
        console.log("Websocket closed");
        webSocket = null;
    };
    webSocket.onmessage = function (evt) {
        var telemetryData = JSON.parse(evt.data);
        render(telemetryData);
    };
    webSocket.onerror = function (evt) {
        console.log("Websocket error");
    };
}
