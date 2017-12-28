angular.module('BoilerAdmin').controller('BoilerRuntimeController', function($rootScope, $scope, $http, $location, $timeout, $uibModal, $log, $state, $document, moment, settings) {
    bRuntime = this;

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components

        App.initAjax();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    bRuntime.init = function () {
        // console.error("runtime init");
        $rootScope.boiler = null;
        $rootScope.instants = [];

        bRuntime.hasBoiler = false;
        bRuntime.hasInstants = false;
        bRuntime.hasRuntime = false;

        var svg = d3.select("svg");
        if (svg) {
            // console.error("Find SVG Init & Remove IT");
            svg.remove();
        }
    };

    bRuntime.init();

    bRuntime.fetchBoiler = function () {
        //var url = $location.url();
        var p = $location.search();
        // $log.error('/boiler_list/?boiler=' + p['boiler'], p)
        $http.get('/boiler_list/?boiler=' + p['boiler'])
            .then(function (res) {
                // $scope.parameters = data;
                $log.error('Runtime Boiler Get:', res);
                bRuntime.boiler = res.data[0];
                bRuntime.hasBoiler = true;
                bRuntime.fetchStatus(bRuntime.boiler);
                setTimeout(function () {
                    App.stopPageLoading();
                }, 300);
            });
    };

    bRuntime.fetchStatus = function (boiler) {
        // $log.error("bRuntime.fetchStatus(boiler):", boiler);
        if (!boiler) {
            return;
        }

        if (!$state.includes('runtime')) {
            $log.info("!$state.includes('runtime') BREAK");
            return;
        }

        $http.get('/boiler/state/is_burning/?boiler=' + boiler.Uid)
            .then(function (res) {
                console.log("Fetch BurningStatus Resp:", res.data);
                boiler.isBurning = (res.data.value === "true");
            }, function (err) {
                console.error('Fetch Status Err!', err);
                boiler.isBurning = false;
            })
            .then(function () {
                bRuntime.fetchRuntime(bRuntime.boiler);
            });

        $http.get('/boiler/state/has_subscribed/?boiler=' + boiler.Uid + "&uid=" + $rootScope.currentUser.Uid)
            .then(function (res) {
                console.log("Fetch SubscribeStatus Resp:", res.data);
                boiler.hasSubscribed = (res.data.value === "true");
            }, function (err) {
                console.error('Fetch Status Err!', err);
            });

        if (boiler.Terminal) {
            $http.get('/boiler/state/has_channel_custom/?terminal=' + boiler.Terminal.Uid + "&uid=" + $rootScope.currentUser.Uid)
                .then(function (res) {
                    boiler.hasChannelCustom = res.data['HasCustom'];
                }, function (err) {
                    console.error('Fetch Status Err!', err);
                });
        }

        setTimeout(function () {
            bRuntime.fetchStatus(boiler);
        }, 15000);
    };

    bRuntime.fetchRuntime = function (boiler) {
        //var rtmQ = [1001, 1002, 1003, 1014, 1015, 1016, 1005, 1006, 1017, 1018, 1019, 1021, 1201, 1202];
        var rtmQ = [];

        $rootScope.statusMode = 0;

        //TODO: For Spec
        // console.error("Boiler For Instants:", boiler);
        var ter = boiler.Terminal;
        if (ter) {
            switch (ter.TerminalCode) {
                case 680055:
                    $rootScope.statusMode = 1;
                    break;
                case 680082:
                case 680085:
                case 680096:
                case 680120:
                    $rootScope.statusMode = 2;
                    break;
                case 680064:
                    $rootScope.statusMode = 3;
                    break;
                case 680500:
                case 680053:
                case 680501:
                case 680502:
                    $rootScope.statusMode = 5;
                    break;
            }
        }

        var data = {
            uid: boiler.Uid,
            runtimeQueue: rtmQ,
            limit: 50
        };

        bRuntime.data = { Uid: bRuntime.boiler.Uid };

        $http.post('/boiler_runtime_instants/', data).then(function (res) {
            $log.info("instants Resp:", res);

            boiler.alarmLevel = boiler.isBurning ? 0 : -1;
            boiler.hasSwitchValue = false;

            var instants = [];
            for (var i = 0; i < res.data.length; i++) {
                var d = res.data[i];
                var value;
                var name = d.ParameterName;
                var alarmLevel = -1;

                /*
                if (boiler.isBurning) {
                    value = d.Value;
                    alarmLevel = d.AlarmLevel;
                } else {
                    value = "-";
                }
                */
                value = d.Value;
                alarmLevel = d.AlarmLevel;

                if (alarmLevel > boiler.alarmLevel) {
                    boiler.alarmLevel = alarmLevel;
                }

                var label = "";
                switch (alarmLevel) {
                    case -1:
                        label = "未测定";
                        break;
                    case 0:
                        label = "正常";
                        break;
                    case 1:
                        //label = rtm.Alarm.AlarmNormal > rtm.Alarm.AlarmWarning ? "过低" : "过高";
                        label = "超标";
                        break;
                    case 2:
                        label = "告警";
                        break;
                }

                if (bRuntime.boiler.Form.Id === 205) {
                    switch (d.Parameter) {
                        case 1005:
                            name = "回水温度";
                            break;
                        case 1006:
                            name = "出水温度";
                            break;
                    }
                }

                if (d.ParameterCategory === 11) {
                    boiler.hasSwitchValue = true;
                }

                instants.push({
                    id: d.Parameter,
                    name: name,
                    category: d.ParameterCategory,
                    value: value,
                    unit: d.Unit,
                    alarmLevel: alarmLevel,
                    alarmDesc: label,
                    date: new Date(d.UpdatedDate)
                });

                bRuntime.instants = instants;
            }

            $rootScope.boiler = bRuntime.boiler;
            $rootScope.instants = bRuntime.instants;

            $log.info("Boiler Inst:", bRuntime.instants);
            bRuntime.hasInstants = true;
        }, function (err) {
            //alert('Fetch Err!' + err.status + " | " + err.data);
        }).then(function () {

        });


        $http.post('/boiler_runtime_list/', data).then(function (res) {
            console.info("Runtime Resp:", res);
            // alert("Boiler Put Detail Res," + res.status + res.data + "|" + Object.keys(res.data))
            // bRuntime.boiler.Runtimes = res.data.Runtimes;
            // bRuntime.boiler.Parameters = res.data.Parameters;
            // bRuntime.boiler.Rules = res.data.Rules;

            if (res.data.Parameters) {
                for (var i = 0; i < res.data.Parameters.length; i++) {
                    var param = res.data.Parameters[i];
                    var pid = param.Id;

                    bRuntime.data[pid] = res.data.Runtimes[i];
                }
            }
            //console.error("Runtime Data: ", bRuntime.data);
            /*
            for (var i = 0; i < bRuntime.boiler.Parameters.length; i++) {
                var rtm, value;
                var aParam = bRuntime.boiler.Parameters[i];
                //var name = aParam.Name.replace(/([A-Z]+)([1-9])/, "$1<sub>$2</sub>");
                var name = aParam.Name;
                var unit = aParam.Unit;
                var id = aParam.Id;
                var alarmLevel = -1;

                if (bRuntime.boiler.Runtimes[i] && bRuntime.boiler.Runtimes[i].length > 0) {
                    rtm = bRuntime.boiler.Runtimes[i][0];
                    value = (rtm.value * aParam.Scale).toFixed(aParam.Fix);
                    if (rtm.alarm_level && rtm.alarm_level.length > 0) {
                        alarmLevel = parseInt(rtm.alarm_level);
                    } else {
                        alarmLevel = 0;
                    }
                } else {
                    value = "-";
                }

                var label = "";
                switch (alarmLevel) {
                    case -1:
                        label = "未测定";
                        break;
                    case 0:
                        label = "正常";
                        break;
                    case 1:
                        //label = rtm.Alarm.AlarmNormal > rtm.Alarm.AlarmWarning ? "过低" : "过高";
                        label = "超标";
                        break;
                    case 2:
                        label = "告警";
                        break;
                }

                bRuntime.boiler.data.push({
                    id: id,
                    name: name,
                    unit: unit,
                    value: value,
                    alarmLevel: alarmLevel,
                    label: label,
                    date: rtm ? rtm.created_date : ""
                });
            }
            */
            $rootScope.boilerRuntime = bRuntime.data;
            console.info("Boiler Runtime Data:", bRuntime.boiler);

            bRuntime.hasRuntime = true;
        }).then(function () {
            //TODO: MAYBE HAS ISSUE;
            // bRuntime.initCharts();
            bRuntime.fetchDaily();
        });
    };

    bRuntime.fetchDaily = function () {
        var p = $location.search();
        var limit = 30;
        $http.post('/boiler_runtime_daily/', {
            uid: p['boiler'],
            limit: limit
        }).then(function (res) {
            console.warn("Runtime Flows Resp:", res);
            var pa = res.data.Parameter;

            bRuntime.daily = [];
            for (var i = 0; i < limit; i++) {
                var flow = res.data.Flows && i < res.data.Flows.length && res.data.Flows[i] ?
                    res.data.Flows[i].Value : 0;
                var heat = res.data.Heats && i < res.data.Heats.length && res.data.Heats[i] ?
                    res.data.Heats[i].Value : 0;
                //console.info("Total OBJ:", total);
                var aDay = new Date();
                aDay.setHours(0);
                aDay.setMinutes(0);
                aDay.setSeconds(0);
                aDay.setDate(aDay.getDate() - i);
                var date = res.data.Flows && i < res.data.Flows.length && res.data.Flows[i] ?
                    new Date(res.data.Flows[i].Date) : aDay;

                bRuntime.daily.push({
                    flow: flow.toFixed(2),
                    heat: heat.toFixed(2),
                    date: date
                });
            }

            $rootScope.bRuntime = bRuntime.daily;
            console.info("BoilerData:", bRuntime.daily);

            initChartHeatMonth(bRuntime.daily);
        });
    };

    bRuntime.setSubscribe = function (boiler) {
        $http.post('/boiler/state/set_subscribe/', {
            uid: boiler.Uid,
            state: bRuntime.boiler.hasSubscribed
        }).then(function (res) {
            console.info("Set Subscribe Resp:", res);
        });
    };

    bRuntime.initCharts = function (boiler) {
        console.info("Runtime initCharts!");

        // initChartSmokeComponentsAm(boiler);
        // initChartSteamAm(boiler);
        // initChartTemperatureAm(boiler);
        // initChartExcessAir(boiler);
        // initChartHeat(boiler);
        initChartHeatMonth(boiler);
    };
});

var bRuntime;

function boiler_module_height() {
    var mo = document.getElementById('module_height');
    if (!!window.ActiveXObject || "ActiveXObject" in window) {
        //判断是否为IE
        console.warn("IsIE");
        if (mo) {
            //已获取到module_height元素,等比缩放1.5x
            mo.setAttribute('height', document.documentElement.clientWidth * 0.50);
        }
        else {
            var MutationObserver = window.MutationObserver ||
                window.WebKitMutationObserver ||
                window.MozMutationObserver;
            var mutationObserverSupport = !!MutationObserver;
            if (mutationObserverSupport) {
                //判断是否支持mutationObserver
                document.getElementById('boiler_module').addEventListener("DOMSubtreeModified", function () {
                    boiler_module_height();
                    console.log('DOMNodeInserted');
                }, false);
            }
            else {
                //不支持mutationObserver使用DOMNodeInserted触发器
                $("#boiler_module").bind('DOMNodeInserted', function (e) {
                    boiler_module_height();
                });
            }
        }
    }
    else {
        console.info("Not Is IE")
    }
}

//window.onresize = function () {
//    boiler_module_height();
//};