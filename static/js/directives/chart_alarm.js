/**
 * Created by JeremiahYan on 2017/5/8.
 */
boilerAdmin.directive('chartAlarm', function () {
    return {
        restrict: 'E',
        templateUrl: "/directives/chart_alarm.html",
        link: function(scope, element, attrs) {
            initChartAlarm(scope.alarm);
        }
    };
});

//2017-01-03T14:22:54+08:00
var initChartAlarm = function(alarm) {
    console.warn("initChartAlarm", alarm);
    if (!alarm) {
        console.warn("Boiler Alarm Data IS NULL!");
        return;
    }
    if (typeof(AmCharts) === 'undefined' || $('#chart_alarm').size() === 0) {
        console.warn("There IS NO #chart_alarm");
        return;
    }

    var chartData = [];

    var pName = alarm['Parameter__Name'];
    var unit = alarm['Parameter__Unit'];
    var scale = alarm['Parameter__Scale'];
    var fix = alarm['Parameter__Fix'];
    var normalValue = alarm['TriggerRule__Normal'];
    var warningValue = alarm['TriggerRule__Warning'];

    for (var i = 0; i < alarm.runtime.length; i++) {
        var rtm = alarm.runtime[i];
        var value = (rtm.Value * scale).toFixed(fix);

        var d = {
            num: i,
            date: new Date(rtm.CreatedDate),
            value: value
        };
        chartData.push(d);
    }

    var lowColor = warningValue > normalValue ? "#0d8ecf" : "#f0868e";
    var highColor = warningValue > normalValue ? "#f0868e" : "#0d8ecf";

    console.info("ChartAlarm Data:", chartData);

    var chart = AmCharts.makeChart("chart_alarm", {
        type: "serial",
        theme: "light",
        fontSize: 11,
        color: "#6c7b88",
        language: "zh",
        marginTop: 8,
        marginBottom: 0,
        marginLeft: 0,
        marginRight: 0,
        dataProvider: chartData,
        // valueAxes: [{
        //     axisAlpha: 0,
        //     position: "left"
        // }],
        valueAxes: [{
            id: "a1",
            //title: "Test",
            //position: "right",
            gridAlpha: 0.04,
            axisAlpha: 0.2,
            unit: unit,
        }],
        graphs: [{
            id:"g1",
            valueField: "value",
            balloonText: pName + ":[[value]]" + unit,//"[[category]]<br><b><span style='font-size:14px;'>[[value]]</span></b>",
            //bullet: "round",
            //bulletSize: 2,
            lineColor: highColor,
            negativeLineColor: lowColor,
            lineThickness: 1,
            negativeBase: warningValue,
            fillAlphas: 0.3,
            lineAlpha: 0.6,
            type: "smoothedLine"
        }],
        /*
        "chartScrollbar": {
            "graph":"g1",
            "gridAlpha":0,
            "color":"#888888",
            "scrollbarHeight":55,
            "backgroundAlpha":0,
            "selectedBackgroundAlpha":0.1,
            "selectedBackgroundColor":"#888888",
            "graphFillAlpha":0,
            "autoGridCount":true,
            "selectedGraphFillAlpha":0,
            "graphLineAlpha":0.2,
            "graphLineColor":"#c2c2c2",
            "selectedGraphLineColor":"#888888",
            "selectedGraphLineAlpha":1

        },
        */
        chartCursor: {
            cursorAlpha: 0,
            categoryBalloonDateFormat: "JJ:NN",
            categoryBalloonColor: "#e26a6a",
            categoryBalloonAlpha: 0.8,

            valueLineEnabled: true,
            valueLineBalloonEnabled: true,
            valueLineAlpha: 0.3,
            bulletsEnabled: true,
            bulletSize: 8,
            fullWidth:true
        },
        dataDateFormat: "YYYY-MM-DD JJ:NN:SS",
        categoryField: "date",
        categoryAxis: {
            minPeriod: "mm",
            parseDates: true,
            equalSpacing: true,
            axisAlpha: 0.2,
            gridAlpha: 0.04,
            // minorGridAlpha: 0.1,
            // minorGridEnabled: true
        }
    });
};
