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
                if($state.includes("runtime.stats")){
                    bRuntime.initList();
                }
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

        if ($location.search()['boiler'] !== boiler.Uid) {
            $log.info("Boiler Uid Changed!");
            return;
        }

        $http.get("/terminal_boiler_status/?boiler="+ boiler.Uid)
            .then(function (res) {
                boiler.isOnline = res.data.term_status;
                boiler.isBurning = res.data.boiler_status && boiler.isOnline;

            },function (err) {
                boiler.isOnline = false;
                boiler.isBurning = false;
            }).then(function () {
            $rootScope.isBoilerOnline = boiler.isOnline;
            $rootScope.isBoilerBurning = boiler.isBurning;
            bRuntime.fetchRuntime(bRuntime.boiler);
        });

        /*$http.get('/boiler/state/is_online/?boiler=' + boiler.Uid)
            .then(function (res) {
                console.info("Fetch OnlineStatus Resp:", res.data);
                boiler.isOnline = res.data;
            }, function (err) {
                console.error('Fetch Status Err!', err);
                boiler.isOnline = false;
            })
            .then(function () {
                $rootScope.isBoilerOnline = boiler.isOnline;
                $http.get('/boiler/state/is_burning/?boiler=' + boiler.Uid)
                    .then(function (res) {
                        console.info("Fetch BurningStatus Resp:", res.data);
                        boiler.isBurning = (res.data.value === "true");
                    }, function (err) {
                        console.error('Fetch Status Err!', err);
                        boiler.isBurning = false;
                    })
                    .then(function () {
                        $rootScope.isBoilerBurning = boiler.isBurning;
                        // boiler.alarmLevel = ( boiler.isOnline && boiler.isBurning ) ? 0 : -1;
                        bRuntime.fetchRuntime(bRuntime.boiler);
                    });
            });*/



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
        }, 5000);
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


            boiler.hasSwitchValue = false;
            boiler.hasRangeValue = false;
            boiler.alarmLevel = ( boiler.isOnline && boiler.isBurning ) ? 0 : -1;

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
                if (d.ParameterCategory === 11) {
                    d.SwitchFlag = d.SwitchStatus;
                    d.AlarmLevel = 0;
                }

                if (boiler.isBurning && boiler.isOnline) {
                    alarmLevel = d.AlarmLevel;
                }


                if (alarmLevel > boiler.alarmLevel) {
                    boiler.alarmLevel = alarmLevel;
                }
                // console.log(d.AlarmLevel);

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
                /*if (bRuntime.boiler.Form.Id === 205) {
                    switch (d.Parameter) {
                        case 1005:
                            name = "回水温度";
                            break;
                        case 1006:
                            name = "出水温度";
                            break;
                    }
                }*/

                if (d.ParameterCategory === 11) {
                    boiler.hasSwitchValue = true;
                }

                if (d.ParameterCategory === 13) {
                    boiler.hasRangeValue = true;
                }

                var inst = {
                    id: d.Parameter,
                    name: name,
                    category: d.ParameterCategory,
                    value: d.ParameterCategory !== 13 ? value : d.Remark,
                    switchFlag: d.SwitchFlag,
                    unit: d.Unit,
                    alarmLevel: alarmLevel,
                    alarmDesc: label,
                    date: new Date(d.UpdatedDate),
                    remark:d.Remark
                };

                instants.push(inst);

                if (!bRuntime.currentDate || inst.date > bRuntime.currentDate) {
                    bRuntime.currentDate = inst.date;
                }
            }

            bRuntime.instants = instants;

            $rootScope.boiler = bRuntime.boiler;
            $rootScope.instants = bRuntime.instants;

            $log.info("Boiler Inst:", bRuntime.instants);
            bRuntime.hasInstants = true;
        }, function (err) {
            //alert('Fetch Err!' + err.status + " | " + err.data);
        }).then(function () {

        });



    };

    bRuntime.initList = function () {
        var rtmQ = [];
        var p = $location.search();

        if(!bRuntime.boiler){
            return;
        }
        var data = {
            uid: bRuntime.boiler.Uid,
            runtimeQueue: rtmQ,
            limit: 50
        };

        bRuntime.data = { Uid: bRuntime.boiler.Uid };

        console.log("bRuntime.boiler:",bRuntime.boiler);
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


            if($state.includes("runtime.stats")){
                setTimeout(function () {
                    bRuntime.initList();
                },15000);
            }


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

    $scope.kai=true;
    $scope.valueData ="正常";
    $scope.kaiguan=function () {
        $scope.kai = !$scope.kai;
        if($scope.kai){
            $scope.valueData ="正常";
        }else{
            $scope.valueData ="调试"
        }
    }


    //---------锅炉重启------------
    /*$scope.terRes =[
        {
            value:2,
            name:"启动锅炉"
        },
        {
            value:1,
            name:"关闭锅炉"
        },
        {
            value:1,
            name:"重置锅炉"
        }
    ];*/
    bRuntime.initControl = function () {
        var p = $location.search();
        $http.post("/issued_boiler_status",{
            uid: p['boiler']
        }).then(function (res) {
            $scope.isRemoteControl = res.data;
        },function (err) {

        });
    };

    $scope.boilerRestart = function(value){
        var ter = bRuntime.boiler.TerminalsCombined[0];
        console.log(ter);
        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/issued_boiler",{
            boiler_id:bRuntime.boiler.Uid,
            terminal_code:ter.TerminalCode.toString(),
            value:parseInt(value)
        }).then(function (res) {
            App.stopPageLoading();
            swal({
                title: "成功",
                text: res.data,
                type: "success"
            });
        },function (err) {
            App.stopPageLoading();
            swal({
                title: "失败",
                text: err.data ,
                type: "error"
            });
        });
    }



//    ---------参数设置---------

    $scope.dataConfig = function () {
        var modalInstance = $uibModal.open({
            templateUrl: '/directives/modal/runtime_config.html',
            controller: 'runtimeConfigCtrl',
            // backdrop: "static",
            size: "",
            windowClass: 'zindex_sub',
            // resolve: {
            //     items1: function () {
            //         return $scope.items;
            //     }
            // }
        });

        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
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


angular.module('BoilerAdmin').controller("statusModule", function($scope,$rootScope) {
    $scope.boiler = [];
    $scope.valueLabels = {};
    $scope.statusLabels = {};
    $scope.switchLabels = {};
    $rootScope.$watch('boiler', function () {
        // console.error("$rootScope.$watch('boiler')", $rootScope.boiler, bModule.boiler);
        if (!$rootScope.boiler || ($scope.boiler === $rootScope.boiler)){
            return;
        }
        $scope.boiler = $rootScope.boiler;
        $scope.initStatus();
    });

    $rootScope.$watch('instants', function () {
        // console.error("$rootScope.$watch('instants')", bModule.instants);
        if (!$rootScope.instants) {
            return;
        }
        $scope.instants = $rootScope.instants;

        $scope.updateStatusLabels();
        // $scope.updateLabels();
    });

    $rootScope.$watch('isBoilerBurning', function () {
        // console.error("$rootScope.$watch('isBoilerBurning')", $rootScope.isBoilerBurning);
        $scope.isBoilerBurning = $rootScope.isBoilerBurning;

        $scope.updateStatusLabels();
    });

    var moduleStatus = d3.select("#status_1");
    var svgContainer = moduleStatus.append("svg");
    $scope.initStatus = function () {
        // if (!$scope.boiler ||
        //     !$rootScope.instants ||
        //     $rootScope.instants.length <= 0 ||
        //     $scope.instants === $rootScope.instants) {
        //     return;
        // }
        $scope.instants = $rootScope.instants;

        var isTerminalConnected = $scope.boiler && $scope.boiler.isOnline;
        var statData = [
            [{
                id: 1,
                name: "终端状态",
                text: isTerminalConnected ? "已连接" : "未连接",
                type: "status",
                value: !!isTerminalConnected
            },
                {
                    id: 2,
                    name: "运行状态",
                    text:($scope.boiler.isOnline && $scope.boiler.isBurning) ? "正在运行" : "未运行",
                    type: "status",
                    value: ($scope.boiler.isOnline && $scope.boiler.isBurning)
                },
                {
                    id: 3,
                    name: "告警状态",
                    text: "",
                    type: "status",
                    value: $scope.boiler.alarmLevel
                }
            ],
//		[{
//				id: 0,
//				name: "热效率(正平衡)"
//			},
//			{
//				id: 1201,
//				name: "热效率(反平衡)"
//			}
//		]

        ];
        // console.log($scope.boiler.alarmLevel);
        var statOptions = {
            align: "left", //"left", "justify"
            baseWidth: 82,
            height: 40,
            gap: 8,
            baseX: 0,
            baseY: 0
        };;
        statOptions.align = "justify";
        renderStatusModule(statData, statOptions);
    };



    var renderStatusModule = function(data, options) {

        // console.log("ready to renderStatusModule", data, options);
        var align = options.align;

        var baseWidth = options.baseWidth;
        var height = options.height;
        var gap = options.gap;
        var fontSize = Math.round(baseWidth / 7);

        var baseX = options.baseX;
        var baseY = options.baseY;

        var statusModule = svgContainer;

        if(!statusModule) {
            $log.warn("There IS NO " + id + "!");
            return;
        }

        var maxRowLength = 0;
        for(var row = 0; row < data.length; row++) {
            if(data[row].length > maxRowLength) {
                maxRowLength = data[row].length;
            }
        }

        for(var row = 0; row < data.length; row++) {
            var rowData = data[row];
            for(var col = 0; col < rowData.length; col++) {
                var width, cx, cy;
                cy = baseY + (height + gap) * row;

                switch(align) {
                    case "left":
                        width = baseWidth;
                        cx = baseX + (width + gap) * col;
                        break;
                    case "right":
                        width = baseWidth;
                        cx = baseX + (width + gap) * (maxRowLength - rowData.length) + (width + gap) * col;
                        break;
                    case "justify":
                        width = (baseWidth * maxRowLength + gap * (maxRowLength - rowData.length)) / rowData.length;
                        cx = baseX + (width + gap) * col;
                        break;
                    default:
                        width = baseWidth;
                        cx = baseX + (width + gap) * col;
                        break;
                }

                var d = rowData[col];


                var barColor = "#4c87b9";
                var text = (d.type === "switch" ? "" : (d.type === "status" ? d.text : "未测定"));
                var textColor = d.type === "status" ? "#fff" : "#aaa";


                if ($scope.boiler.isBurning && d.type !== "status" && d.id > 0) {
                    for (var i = 0; i < $scope.instants.length; i++) {
                        var ins = $scope.instants[i];
                        if (d.id == ins.id && ins.value != "-") {
                            barColor = "#4c87b9";
                            text = ins.value + ins.unit;
                            textColor = "#80898e";
                            break;
                        }
                    }
                }



                //Bar Drawing
                statusModule.append("rect")
                    .attr("x", cx)
                    .attr("y", cy)
                    .attr("width", width)
                    .attr("height", height)
                    //.attr("rx", 6)
                    .style("fill", "none")
                    .style("stroke", barColor)
                    .style("stroke-width", "1");
                statusModule.append("rect")
                    .attr("x", cx)
                    .attr("y", cy)
                    .attr("width", width)
                    .attr("height", height / 2)
                    .style("fill", barColor);


                if(d.type === "status") {
                    //StatusColor Drawing
                    var bgColor = "#32c5d2";
                    if(typeof d.value === "boolean") {
                        bgColor = d.value ? "#32c5d2" : "#e7505a";
                    } else if(typeof d.value === "number") {
                        switch(d.value) {
                            case -1:
                                bgColor = "#cfdae1";
                                break;
                            case 0:
                                bgColor = "#32c5d2";
                                break;
                            case 1:
                                bgColor = "#f3c200";
                                break;
                            case 2:
                                bgColor = "#e7505a";
                                break;
                        }
                    }

                    var statusLabel = statusModule.append("rect")
                        .attr("x", cx + 4)
                        .attr("y", cy + height / 2 + 4)
                        .attr("width", width - 8)
                        .attr("height", height / 2 - 8)
                        .attr("rx", 6)
                        .attr("ry", 6)
                        .style("fill", bgColor);

                    $scope.statusLabels[d.id] = {};
                    $scope.statusLabels[d.id].label = statusLabel;

                    // statusModule.append("rect")
                    //     .attr("x", cx + 4)
                    //     .attr("y", cy + height / 2 + 4)
                    //     .attr("width", width - 8)
                    //     .attr("height", height / 2 - 8)
                    //     .attr("rx", 6)
                    //     .attr("ry", 6)
                    //     .style("fill", bgColor);
                }

                //Label Drawing
                statusModule.append("text")
                    .attr("x", cx + width / 2)
                    .attr("y", cy + fontSize / 2 + 2)
                    .attr("dy", fontSize / 2)
                    .attr("text-anchor", "middle")
                    .text(d.name)
                    .style("font-size", fontSize + "px")
                    //.style("font-weight", "bold")
                    .style("fill", "#fff")
                    .style("stroke-width", "0px");

                //Text Drawing
                // statusModule.append("text")
                //     .attr("x", cx + width / 2)
                //     .attr("y", cy + height / 2 + fontSize / 2 + 2)
                //     .attr("dy", fontSize / 2)
                //     .attr("text-anchor", "middle")
                //     .text(text)
                //     .style("font-size", fontSize - 2 + "px")
                //     //.style("font-weight", "bold")
                //     .style("fill", textColor)
                //     .style("stroke-width", "0px");

                //Text Drawing

                var valueLabel = statusModule.append("text")
                    .attr("x", cx + width / 2)
                    .attr("x", cx + width / 2)
                    .attr("y", cy + height / 2 + fontSize / 2 + 2)
                    .attr("dy", fontSize / 2)
                    .attr("text-anchor", "middle")
                    .text(text)
                    .style("font-size", fontSize - 2 + "px")
                    .style("fill", textColor)
                    .style("stroke-width", "0px");

                $scope.statusLabels[d.id].value = valueLabel;



            }
        }
    };



    $scope.updateStatusLabels = function () {
        if (!$scope.boiler) {
            if (!$rootScope.boiler) {
                return
            }

            $scope.boiler = $rootScope.boiler;
        }

        var isTerminalConnected = $scope.boiler && $scope.boiler.isOnline;

        for (var i in $scope.statusLabels) {
            var statusLabel = $scope.statusLabels[i];
            var text = "";
            var bgColor = "#32c5d2";

            switch (parseInt(i, 10)) {
                case 1:
                    text = isTerminalConnected ? "已连接" : "未连接";
                    bgColor = isTerminalConnected ? "#32c5d2" : "#bfcad1";
                    break;
                case 2:
                    text = ($scope.isBoilerBurning && isTerminalConnected) ? "正在运行" : "未运行";
                    bgColor = ($scope.isBoilerBurning && isTerminalConnected) ? "#32c5d2" : "#e7505a";
                    break;
                case 3:
                    switch ($scope.boiler.alarmLevel) {
                        case -1:
                            bgColor = "#cfdae1";
                            break;
                        case 0:
                            bgColor = "#32c5d2";
                            break;
                        case 1:
                            bgColor = "#f3c200";
                            break;
                        case 2:
                            bgColor = "#e7505a";
                            break;
                    }
                    break;
            }

            // console.error("updateStatusLabels", i, text, bgColor, statusLabel);

            statusLabel.value.text(text);
            statusLabel.label.style("fill", bgColor);
        }
    };



});


//参数设置
angular.module('BoilerAdmin').controller('runtimeConfigCtrl', function ($scope, $uibModalInstance) {

    $scope.items = {
        steamPress : "0.00",
        controlPress :"0.00",
        smokeTemperature:"0.0",
        pipeTemperature:"0.0"
    };

    $scope.ok = function () {
        $uibModalInstance.close();
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});