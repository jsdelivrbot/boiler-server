/**
 * Created by JeremiahYan on 2017/1/8.
 */
boilerAdmin.directive('chartSteam', function () {
    return {
        restrict: 'E',
        templateUrl: "/directives/chart_steam.html",
        controller: "ChartSteamController",
        controllerAs: "chartSteam"
    };
}).controller("ChartSteamController", function ($rootScope, $scope, $http) {
    var bChart = this;
    bChart.range = RUNTIME_RANGE_DEFAULT;
    var domId = "chart_steam";
    var dataField = "steams";

    var pids = {
        1001: "temperature",
        1002: "pressure",
        1003: "flow"
    };

    var keys = [];
    for (var i = 0; i < Object.keys(pids).length; i++) {
        var k = Object.keys(pids)[i];
        keys.push(parseInt(k));
    }

    $rootScope.$watch('boilerRuntime', function () {
        if (bChart.range === RUNTIME_RANGE_DEFAULT) {
            bChart.refreshChart();
        }
    });

    $scope.$watch(dataField, function () {
        // console.error("$scope.$watch('datasource', function ()", $scope.datasource);
        if (!$scope[dataField]) {
            console.warn("There is no $scope(" + dataField + ")!!");
            return;
        }

        var chartData = [];
        for (var i = 0; i < Object.keys(pids).length; i++) {
            var id = Object.keys(pids)[i];
            var field = pids[id];
            if (!$scope[dataField] || !$scope[dataField][id]) {
                break;
            }
            var rtms = $scope[dataField][id];
            if (!rtms) {
                continue;
            }
            for (var j = 0; j < rtms.length; j++) {
                var r = rtms[j];
                var d = {};
                if (j < chartData.length) {
                    d = chartData[chartData.length - 1 - j];
                } else {
                    d = {};
                    d.num = j;
                    d.date = r.CreatedDate;
                    d[field] = 0;
                }

                d[field] = r.Value;

                if (j >= chartData.length) {
                    chartData.unshift(d);
                }
            }
        }

        console.info("ChartData:", chartData);

        for (var i = 0; i < bChart.chart.graphs.length; i++) {
            var g = bChart.chart.graphs[i];
            if (chartData.length > 160) {
                g.bullet = "none";
            } else {
                g.bullet = "round";
            }
        }

        bChart.chart.dataProvider = chartData;
        bChart.chart.write(domId);
        bChart.chart.validateData();
    });

    bChart.initChart = function () {
        var chart = new AmCharts.AmSerialChart();
        chart.fontSize = 11;
        chart.theme = AmCharts.themes.light;
        chart.color = "#6c7b88";
        chart.language = "zh";

        chart.marginTop = 8;
        chart.marginBottom = 0;
        chart.marginLeft = 0;
        chart.marginRight = 0;

        chart.fontFamily = "Open Sans";
        chart.dataDateFormat = 'YYYY-MM-DD HH:NN:SS';

        // chart.addClassNames = true;
        chart.startDuration = 0;

        chart.categoryField = "date";

        var categoryAxis = new AmCharts.CategoryAxis();
        categoryAxis.minPeriod = "mm";
        categoryAxis.parseDates = true;
        categoryAxis.equalSpacing = true;
        categoryAxis.axisAlpha = 0.2;
        categoryAxis.gridAlpha = 0.04;
        // gridCount: 50,
        // gridColor: "#FFFFFF",
        // axisColor: "#555555",
        // dateFormats: [{
        //     period: 'mm',
        //     format: 'HH:NN:SS'
        // }, {
        //     period: 'hh',
        //     format: 'JJ:NN'
        // }, {
        //     period: 'DD',
        //     format: 'MMM HDD'
        // }, {
        //     period: 'WW',
        //     format: 'MMM DD'
        // }]

        chart.valueAxes = [{
            id: "a1",
            //title: "Test",
            position: "right",
            gridAlpha: 0,
            axisAlpha: 0.2,
            unit: "t/h"
        }, {
            id: "a2",
            position: "left",
            gridAlpha: 0,
            axisAlpha: 0.2,
            // labelsEnabled: true,
            unit: "℃"
        }, {
            id: "a3",
            position: "left",
            gridAlpha: 0,
            axisAlpha: 0.2,
            inside: true,
            unit: "MPa"
        }];

        var graphFlow = new AmCharts.AmGraph();

        chart.graphs = [{
            id: "g1",
            valueField: "flow",
            title: "蒸汽流量",
            type: "line",
            bullet: "square",
            //bulletSizeField: "townSize",
            bulletBorderColor: "#02617a",
            // bulletBorderAlpha: 1,
            // bulletBorderThickness: 2,
            bulletSize: 2,
            fillAlphas: 0.4,
            valueAxis: "a1",
            balloonText: "瞬时流量:[[value]] t/h",
            legendValueText: "[[value]]",
            legendPeriodValueText: "总计: [[value.sum]]",
            lineColor: "#89c4f4",
            //alphaField: "alpha"
        }, {
            id: "g2",
            valueField: "temperature",
            //classNameField: "bulletClass",
            title: "蒸汽温度",
            type: "smoothedLine",
            valueAxis: "a2",
            lineColor: "#b0de09",
            lineThickness: 1,
            legendValueText: "[[value]] ℃",
            //descriptionField: "date",
            bullet: "round",
            //bulletSizeField: "townSize",
            bulletColor: "#b0de09",
            // bulletBorderColor: "#a5be4b",
            // bulletBorderAlpha: 1,
            // bulletBorderThickness: 2,
            bulletSize: 2,
            //bulletAlpha: 0.6,
            //labelText: "[[townName2]]",
            labelPosition: "right",
            balloonText: "温度:[[value]] ℃",
            //showBalloon: true,
            //animationPlayed: true,
        }, {
            id: "g3",
            title: "蒸汽压力",
            valueField: "pressure",
            type: "smoothedLine",
            valueAxis: "a3",
            lineAlpha: 0.8,
            lineColor: "#e26a6a",
            balloonText: "[[value]] MPa",
            lineThickness: 1,
            legendValueText: "[[value]] MPa",
            bullet: "round", //"triangleUp",
            // bulletBorderColor: "#e26a6a",
            // bulletBorderThickness: 1,
            // bulletBorderAlpha: 0.8,
            bulletSize: 2,
            // dashLengthField: "dashLength",
            animationPlayed: false
        }];

        chart.chartCursor = {
            //zoomable: false,
            bulletsEnabled: true,
            bulletSize: 6,
            categoryBalloonDateFormat: "MMM DD JJ:NN",
            cursorAlpha: 0,
            categoryBalloonColor: "#e26a6a",
            categoryBalloonAlpha: 0.8,
            fullWidth: true
            //valueBalloonsEnabled: false
        };

        chart.categoryAxis = categoryAxis;

        bChart.chart = chart;
    };

    bChart.refreshChart = function (range) {
        // console.error("initChartSteamAm", $rootScope.boilerRuntime);
        if ((!range || range === RUNTIME_RANGE_DEFAULT) && !$rootScope.boilerRuntime) {
            console.warn("ChartSteamAm BoilerRuntimeData IS NULL!");
            return;
        }
        if (typeof(AmCharts) === 'undefined' || $('#' + domId).size() === 0) {
            console.warn("There IS NO #chart_steam");
            return;
        }

        var since;

        if (bChart.range === range && range != RUNTIME_RANGE_DEFAULT) {
            return;
        }

        bChart.range = range;
        switch (range) {
            case RUNTIME_RANGE_TODAY:
            case RUNTIME_RANGE_THERE_DAY:
            case RUNTIME_RANGE_WEEK:
                var postData = {
                    uid: $rootScope.boilerRuntime.Uid,
                    runtimeQueue: keys,
                    range: range
                };
                if (since && typeof since === 'object') {
                    postData.since = since;
                }
                Ladda.create(document.getElementById('chartSteam' + range)).start();
                $http.post('/boiler_runtime_list/', postData).then(function (res) {
                    console.warn("Ranged:", range, "Runtime Resp:", res);

                    var datasource = { Uid: $rootScope.boilerRuntime.Uid };

                    for (var i = 0; i < res.data.Parameters.length; i++) {
                        var param = res.data.Parameters[i];
                        var pid = param.Id;

                        datasource[pid] = res.data.Runtimes[i];
                    }

                    $scope[dataField] = datasource;
                    Ladda.create(document.getElementById('chartSteam' + range)).stop();
                });
                break;
            case RUNTIME_RANGE_DEFAULT:
            default:
                $scope[dataField] = $rootScope.boilerRuntime;
                break;
        }
    };

    bChart.initChart();
});
