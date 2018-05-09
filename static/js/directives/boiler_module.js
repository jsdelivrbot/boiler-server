/**
 * Created by JeremiahYan on 2017/4/10.
 */
boilerAdmin.directive('boilerModule', function () {
    return {
        restrict: 'E',
        templateUrl: "/directives/boiler_module.html",
        controller: "BoilerModuleController",
        controllerAs: "module"
        //bindToController: true
    };
}).controller("BoilerModuleController", function ($rootScope, $scope, $http, $location, $timeout, $log, $document, moment) {
    bModule = this;

    bModule.valueLabels = {};
    bModule.statusLabels = {};
    bModule.switchLabels = {};

    var copy = function (obj) {
        var aObj = {};

        for (var i = 0; i < Object.keys(obj).length; i++) {
            var key = Object.keys(obj)[i];
            var value = obj[key];
            aObj[key] = ((typeof value) === 'object' ? copy(value) : value);
        }

        return aObj;
    };

    var moduleOptionsDef = {
        align: "left",  //"left", "justify"

        baseWidth: 110,
        height: 54,
        gap: 8,

        baseX: 0,
        baseY: 0
    };

    $rootScope.$watch('boiler', function () {
        // console.error("$rootScope.$watch('boiler')", $rootScope.boiler, bModule.boiler);
        if (!$rootScope.boiler || bModule.boiler === $rootScope.boiler) {
            return;
        }

        bModule.boiler = $rootScope.boiler;

        bModule.initModule();
    });

    $rootScope.$watch('instants', function () {
        // console.error("$rootScope.$watch('instants')", bModule.instants);
        if (!$rootScope.instants) {
            return;
        }
        bModule.instants = $rootScope.instants;

        bModule.updateStatusLabels();
        bModule.updateLabels();
    });

    $rootScope.$watch("isBoilerOnline",function () {
        bModule.isOnline = $rootScope.isBoilerOnline ;

        bModule.renderAnimations();
        bModule.updateStatusLabels();
    });

    $rootScope.$watch('isBoilerBurning', function () {
        // console.error("$rootScope.$watch('isBoilerBurning')", $rootScope.isBoilerBurning);
        bModule.isBoilerBurning = $rootScope.isBoilerBurning ;

        bModule.renderAnimations();
        bModule.updateStatusLabels();
    });

    bModule.initModule = function () {
        if (!bModule.boiler ||
            !$rootScope.instants ||
            $rootScope.instants.length <= 0 ||
            bModule.instants === $rootScope.instants) {
            return;
        }
        bModule.instants = $rootScope.instants;

        console.error("Runtime initModule!", bModule.instants);
        var module = d3.select("#boiler_module");
        if (!module) {
            console.warn("There IS NO #boiler_module!");
            return;
        }

        // console.error("Runtime drawModule!", runtime.instants);
        // if (hasInitModule) {
        //     console.error("Already initModule!");
        //     return;
        // }
        // hasInitModule = true;

        var svgName = "/img/boiler_coal_double.svg";
        bModule.moduleId = BOILE_MODULE_COAL_DOUBLE;
        switch (bModule.boiler.Fuel.Type.Id) {
            case 2:
                bModule.moduleId = BOILE_MODULE_OIL;
                break;
            case 3:
                bModule.moduleId = BOILE_MODULE_GAS;
                break;
            default:
                bModule.moduleId = BOILE_MODULE_COAL_WATER;
                break;
        }

        if (bModule.boiler.Medium.Id === 2 ||
            bModule.boiler.Form.Id === 205) {
            bModule.moduleId = BOILE_MODULE_WATER;
        }

        if ((bModule.boiler.Fuel.Type.Id === 1 || bModule.boiler.Fuel.Type.Id === 4) && bModule.boiler.Medium.Id === 2) {
            bModule.moduleId = BOILE_MODULE_COAL_WATER;
        }

        if (bModule.boiler.Form.Id === 1003) {
            bModule.moduleId = BOILE_MODULE_HEAT_WATER_SYSTEM;
        }
        if (bModule.boiler.Form.Id === 1004) {
            bModule.moduleId = BOILE_MODULE_LV;
        }
        if (bModule.boiler.Form.Id === 1005) {
            bModule.moduleId = BOILE_MODULE_IRON;
        }

        if (bModule.boiler.Form.Id === 1006) {
            bModule.moduleId = BOILE_MODULE_HEAT_STEAM_SYSTEM;
        }


        switch (bModule.moduleId) {
            case BOILE_MODULE_COAL_DOUBLE:
                svgName = "/img/boiler_coal_double.svg";
                break;
            case BOILE_MODULE_OIL:
                svgName = "/img/boiler_oil.svg";
                break;
            case BOILE_MODULE_GAS:
                svgName = "/img/boiler_gas.svg";
                break;
            case BOILE_MODULE_WATER:
                svgName = "/img/boiler_water.svg";
                break;
            case BOILE_MODULE_COAL_WATER:
                svgName = "/img/boiler_coal_water.svg";
                break;
            case BOILE_MODULE_HEAT_WATER_SYSTEM:
                svgName = "/img/electricity.svg";
                break;
            case BOILE_MODULE_HEAT_STEAM_SYSTEM:
                svgName = "/img/electricity2.svg";
                break;
            case BOILE_MODULE_LV:
                svgName = "/img/zhulv.svg";
                break;
            case BOILE_MODULE_IRON:
                svgName = "/img/zhutie.svg";
                break;
            default:
                svgName = "/img/boiler_coal_double.svg";
                break;
        }

        if  (bModule.boiler.TerminalCode == 680071) {
            svgName = "/img/boiler_gas_temp.svg";
        }

        // svgName = "/img/boiler_module_test.svg";

        d3.xml(svgName).mimeType("image/svg+xml").get(function(error, xml) {
            if (error) throw error;
            var svgNode = xml.getElementsByTagName("svg")[0];
            if (!svgNode || !module.node()) {
                console.warn("There IS NO SVG Container!");
                return;
            }

            var svg = module.select("svg");
            if (svg) {
                console.warn("Find SVG & Remove IT");
                svg.remove();
            }

            module.node().appendChild(svgNode);
            bModule.svg = module.select("svg");

            bModule.gauge = bModule.svg.select("#gauge_container");
            bModule.dash = bModule.svg.select("#dash_container");

            /*
             var min = 0;
             var max = 100;

             var config = {
             size: 80,
             label: "Gauge",
             x: 20,
             y: 220,
             min: undefined != min ? min : 0,
             max: undefined != max ? max : 100,
             minorTicks: 5
             };

             var range = config.max - config.min;
             config.yellowZones = [{from: config.min + range * 0.75, to: config.min + range * 0.9}];
             config.redZones = [{from: config.min + range * 0.9, to: config.max}];

             gauges[name] = new Gauge("gauge_container", config);
             gauges[name].render();
             */

            var isTerminalConnected = (bModule.boiler.Terminal && bModule.boiler.Terminal.IsOnline) || bModule.boiler.isBurning;
            var statData = [
                [
                    {id: 1, name: "终端状态", text: isTerminalConnected ? "已连接" : "未连接", type: "status", value: !!isTerminalConnected},
                    {id: 2, name: "燃烧状态", text: bModule.boiler.isBurning ? "已点燃" : "未点燃", type: "status", value: bModule.boiler.isBurning},
                    {id: 3, name: "告警状态", text: "", type: "status", value: bModule.boiler.alarmLevel}
                ]/*,
                [
                    {id: 0, name: "热效率(正平衡)"},
                    {id: 1201, name: "热效率(反平衡)"}
                ]*/
                // "蒸汽超压", "环境温度",
                // "运行时间(每天)", "运行时间(累积)"
            ];

            var statOptions = copy(moduleOptionsDef);
            statOptions.align = "justify";

            var insData = [];
            var insG1Data = [];
            var insG2Data = [];

            var insG0Row = 4;
            var insG0Col = 3;

            var insG1Row = 3;
            var insG1Col = 3;

            var insG2Row = 1;
            var insG2Col = 6;

            var insG0X = 0;
            var insG0Y = 160;

            var insG1X = 820;
            var insG1Y = 10;

            var insG2X = 480;
            var insG2Y = 680;

            switch (bModule.moduleId) {
                case BOILE_MODULE_COAL_DOUBLE:
                    insG0Row = 2;
                    insG1Row = 6;
                    insG1Col = 2;
                    break;
                case BOILE_MODULE_COAL_WATER:
                    insG0Row = 3;
                    insG0Y = 120;

                    insG1Col = 2;
                    insG1Row = 6;

                    insG2Row = 3;
                    insG2Col = 11;
                    insG2X = 0;
                    insG2Y = 750;

                    break;
                case BOILE_MODULE_GAS:
                    insG0Row = 6;
                    insG0Col = 3;
                    insG0Y = 140;

                    insG1Row = 2;
                    insG1Col = 8;
                    insG1X = 370;
                    insG1Y = 0;

                    insG2Row = 2;
                    insG2Col = 11;
                    insG2X = 0;
                    insG2Y = 760;

                    break;
                case BOILE_MODULE_WATER:
                    insG0Row = 6;
                    insG0Col = 3;
                    insG0Y = 140;

                    insG1Row = 2;
                    insG1Col = 8;
                    insG1X = 370;
                    insG1Y = 0;

                    insG2Row = 2;
                    insG2Col = 11;
                    insG2X = 0;
                    insG2Y = 760;
                    break;
                case BOILE_MODULE_HEAT_WATER_SYSTEM:
                    insG0Row = 3;
                    insG0Col = 8;
                    insG0X = 370;
                    insG0Y = 0;
                    break;
                default:

                    break;
            }

            var insG0Num = insG0Row * insG0Col;
            var insG1Num = insG1Row * insG1Col + insG0Num;

            for (var i = 0; i < Math.min(bModule.instants.length / insG0Col, insG0Row); i++) {
                var rowData = [];
                for (var j = 0; j < insG0Col; j++) {
                    var idx = i * insG0Col + j;
                    if (idx >= bModule.instants.length) {
                        break;
                    }
                    var ins = bModule.instants[idx];
                    var iData = {
                        id: ins.id,
                        name: ins.name
                    };
                    switch (ins.category) {
                        case 11:
                            iData.type = "switch";
                            break;
                        case 13:
                            iData.type = "range";
                            break;
                    }

                    switch (ins.id) {
                        case 1021:
                            iData.name = "环境温度";
                            break;
                        case 1202:
                            iData.name = "过量空气系数";
                            break;
                    }
                    rowData.push(iData);
                }

                insData.push(rowData);
            }

            for (var i = 0; i < Math.min((bModule.instants.length - insG0Num) / insG1Col, insG1Row); i++) {
                var rowData = [];
                for (var j = 0; j < insG1Col; j++) {
                    var idx = i * insG1Col + j + insG0Num;
                    if (idx >= bModule.instants.length) {
                        break;
                    }
                    var ins = bModule.instants[idx];
                    var iData = {
                        id: ins.id,
                        name: ins.name
                    };
                    switch (ins.category) {
                        case 11:
                            iData.type = "switch";
                            break;
                        case 13:
                            iData.type = "range";
                            break;
                    }

                    switch (ins.id) {
                        case 1021:
                            iData.name = "环境温度";
                            break;
                        case 1202:
                            iData.name = "过量空气系数";
                            break;
                    }
                    rowData.push(iData);
                }

                insG1Data.push(rowData);
            }
            for (var i = 0; i < Math.min((bModule.instants.length - insG1Num) / insG2Col, insG2Row); i++) {
                var rowData = [];
                for (var j = 0; j < insG2Col; j++) {
                    var idx = i * insG2Col + j + insG1Num;
                    if (idx >= bModule.instants.length) {
                        break;
                    }
                    var ins = bModule.instants[idx];
                    var iData = {
                        id: ins.id,
                        name: ins.name
                    };

                    switch (ins.category) {
                        case 11:
                            iData.type = "switch";
                            break;
                        case 13:
                            iData.type = "range";
                            break;
                    }
                    // console.warn("insG2:", ins, i, j, idx);
                    switch (ins.id) {
                        case 1021:
                            iData.name = "环境温度";
                            break;
                        case 1202:
                            iData.name = "过量空气系数";
                            break;
                    }
                    rowData.push(iData);
                }

                insG2Data.push(rowData);
            }

            var insOptions = copy(moduleOptionsDef);
            insOptions.baseX = insG0X;
            insOptions.baseY = insG0Y;

            var insG1Options = copy(moduleOptionsDef);
            insG1Options.baseX = insG1X;
            insG1Options.baseY = insG1Y;

            var insG2Options = copy(moduleOptionsDef);
            insG2Options.baseX = insG2X;
            insG2Options.baseY = insG2Y;

            // console.warn("insData:", insData);

            //TODO: Others (NEED DISCARD)
            var steamData = [];
            var steamOptions = copy(moduleOptionsDef);

            var waterLvData = [
                [{id: 0, type: "status", text: "", name: "高水位", value: -1}],
                [{id: 0, type: "status", text: "", name: "低水位", value: -1}]
            ];
            var waterLvOptions = copy(moduleOptionsDef);

            var waterData = [];
            var waterOptions = copy(moduleOptionsDef);

            var fuelData = [];
            var fuelOptions = copy(moduleOptionsDef);

            var smokeData = [];
            var smokeOptions = copy(moduleOptionsDef);

            var steamBaseData = [];

            if ($rootScope.statusMode === 2) {
                steamBaseData = [
                    {id: 1021, name: "环境温度"},
                    {id: 1080, name: "空气湿度"}
                ];
            } else {
                steamBaseData = [
                    {id: 1001, name: "蒸汽温度"},
                    {id: 1002, name: "蒸汽压力"}
                ];
            }

            /*
            switch (bModule.moduleId) {
                case BOILE_MODULE_OIL:
                case BOILE_MODULE_GAS:
                    var swData = [];
                    var swOptions = copy(moduleOptionsDef);

                    var water1Data = [];
                    var water1Options = copy(moduleOptionsDef);

                    steamData = [
                        steamBaseData,
                        [{id: 1003, name: "蒸汽流量(瞬时)"}],
                        [{id: 0, name: "蒸汽流量(累计)"}]
                    ];

                    waterData = [
                        [
                            {id: 1010, name: "给水流量(瞬时)"},
                            {id: 0, name: "给水量(累计)"}
                        ]
                    ];

                    water1Data = [
                        [{id: 1005, name: "给水温度(冷)"}],
                        [{id: 0, type: "status", text: "", name: "软水硬度", value: -1}]
                    ];

                    smokeData = [
                        [{id: 1014, name: "排烟温度(前)"}],
                        [{id: 1202, name: "过量空气系数"}],
                        [{id: 1016, name: "排烟氧量"}]
                    ];

                    swData = [
                        [
                            {id: 1006, name: "给水温度(热)"},
                            {id: 1015, name: "排烟温度(后)"}
                        ]
                    ];

                    fuelData = [
                        [
                            {id: 0, name: "燃料流量(瞬时)"},
                            {id: 0, name: "燃料量(累计)"}
                        ]
                    ];

                    steamOptions.align = "right";
                    steamOptions.baseX = 320;
                    steamOptions.baseY = 140;

                    waterOptions.baseX = 730;
                    waterOptions.baseY = 630;

                    water1Options.baseX = 1020;
                    water1Options.baseY = 220;

                    waterLvOptions.baseX = 360;
                    waterLvOptions.baseY = 330;

                    smokeOptions.baseX = 730;
                    smokeOptions.baseY = 136;

                    swOptions.baseX = 630;
                    swOptions.baseY = 42;

                    fuelOptions.baseX = 120;
                    fuelOptions.baseY = 600;

                    renderStatusModule("#sw_container", swData, swOptions);
                    renderStatusModule("#water1_container", water1Data, water1Options);

                    break;

                case BOILE_MODULE_WATER:
                    if (bModule.boiler.hasChannelCustom) {
                        break;
                    }
                    var waterOutData = [];
                    var waterOutOptions = copy(moduleOptionsDef);

                    var waterInData = [];
                    var waterInOptions = copy(moduleOptionsDef);

                    waterOutData = [
                        [
                            {id: 1006, name: "出水温度"},
                        ]
                    ];

                    waterInData = [
                        [
                            {id: 1005, name: "回水温度"},
                        ]
                    ];

                    smokeData = [
                        [{id: 1014, name: "排烟温度(前)"}],
                        [{id: 1202, name: "过量空气系数"}],
                        [{id: 1016, name: "排烟氧量"}]
                    ];

                    fuelData = [
                        [
                            {id: 0, name: "燃料流量(瞬时)"},
                            {id: 0, name: "燃料量(累计)"}
                        ]
                    ];

                    waterInOptions.baseX = 460;
                    waterInOptions.baseY = 120;

                    waterOutOptions.baseX = 960;
                    waterOutOptions.baseY = 120;

                    smokeOptions.baseX = 1130;
                    smokeOptions.baseY = 380;

                    fuelOptions.baseX = 180;
                    fuelOptions.baseY = 600;

                    renderStatusModule("#water_in_container", waterInData, waterInOptions);
                    renderStatusModule("#water_out_container", waterOutData, waterOutOptions);

                    break;

                default:
                    if ($rootScope.statusMode === 2) {
                        steamData = [
                            steamBaseData,
                            [
                                {id: 1003, name: "蒸汽流量(瞬时)"},
                                {id: 0, name: "蒸汽流量(累计)"}
                            ]
                        ];
                    } else {
                        steamData = [
                            [
                                {id: 1001, name: "蒸汽温度"},
                                {id: 1002, name: "蒸汽压力"}
                            ],
                            [
                                {id: 1003, name: "蒸汽流量(瞬时)"},
                                {id: 0, name: "蒸汽流量(累计)"}
                            ]
                        ];
                    }

                    waterData = [
                        [{id: 0, type: "status", text: "", name: "软水硬度", value: -1}],

                        [{id: 1006, name: "给水温度(热)"}],
                        [{id: 1005, name: "给水温度(冷)"}],
                        [{id: 1010, name: "给水流量(瞬时)"}],
                        [{id: 0, name: "给水流量(累计)"}]
                    ];

                    smokeData = [
                        [
                            {id: 1202, name: "过量空气系数"},
                            {id: 1016, name: "排烟氧量"},
                            {id: 1014, name: "排烟温度(前)"},
                            {id: 1015, name: "排烟温度(后)"}
                        ]
                    ];

                    fuelData = [
                        [{id: 0, name: "进煤量(瞬时)"}],
                        [{id: 0, name: "进煤量(累计)"}]
                    ];

                    waterLvOptions.baseX = 20;
                    waterLvOptions.baseY = 300;

                    waterOptions.baseX = 850;
                    waterOptions.baseY = 128;

                    smokeOptions.baseX = 640;
                    smokeOptions.baseY = 680;

                    fuelOptions.baseX = 20;
                    fuelOptions.baseY = 440;

                    steamOptions.baseX = 490;
                    steamOptions.baseY = 100;

                    break;
            }

            renderStatusModule("#status_container", statData, statOptions);
            renderStatusModule("#instant_container", insData, insOptions);
            renderStatusModule("#instant_g1_container", insG1Data, insG1Options);
            renderStatusModule("#instant_g2_container", insG2Data, insG2Options);

            if (bModule.moduleId !== BOILE_MODULE_COAL_WATER) {
                renderStatusModule("#steam_container", steamData, steamOptions);
                renderStatusModule("#water_container", waterData, waterOptions);
                renderStatusModule("#water_lv_container", waterLvData, waterLvOptions);
                renderStatusModule("#smoke_container", smokeData, smokeOptions);
                renderStatusModule("#fuel_container", fuelData, fuelOptions);
            }
            */

            bModule.renderAnimations();

        });
    };

    var renderStatusModule = function (id, data, options) {

        $log.info("ready to renderStatusModule", id, data, options);
        var align = options.align;

        var baseWidth = options.baseWidth;
        var height = options.height;
        var gap = options.gap;
        var fontSize = Math.round(baseWidth / 7.5);

        var baseX = options.baseX;
        var baseY = options.baseY;

        var statusModule = bModule.svg.select(id);

        if (!statusModule) {
            $log.warn("There IS NO " + id + "!");
            return;
        }

        var maxRowLength = 0;
        for (var row = 0; row < data.length; row++) {
            if (data[row].length > maxRowLength) {
                maxRowLength = data[row].length;
            }
        }

        for (var row = 0; row < data.length; row++) {
            var rowData = data[row];
            for (var col = 0; col < rowData.length; col++) {
                var width, cx, cy;
                cy = baseY + (height + gap) * row;
                if (bModule.moduleId === BOILE_MODULE_GAS) {
                    cy -= 70;
                }

                switch (align) {
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

                // var barColor = (d.type === "status" && d.value >= 0) || d.type === "switch" ? "#4c87b9" : "#bfcad1";
                var barColor = "#4c87b9";
                var text = (d.type === "switch" ? "" : (d.type === "status" ? d.text : "未测定"));
                var textColor = d.type === "status" ? "#fff" : "#aaa";
                if (bModule.boiler.isBurning &&
                    d.type !== "status" && d.type !== "switch" &&
                    d.id > 0) {
                    for (var i = 0; i < bModule.instants.length; i++) {
                        var ins = bModule.instants[i];
                        if (d.id === ins.id && ins.value !== "-") {
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

                if (d.type === "status") {
                    //StatusColor Drawing
                    var bgColor = "#32c5d2";
                    if (typeof d.value === "boolean") {
                        bgColor = d.value ? "#32c5d2" : "#e7505a";
                    } else if (typeof d.value === "number") {
                        switch (d.value) {
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

                    bModule.statusLabels[d.id] = {};
                    bModule.statusLabels[d.id].label = statusLabel;
                }

                if (d.type === "switch") {
                    //StatusColor Drawing
                    var ins;
                    for (var i = 0; i < bModule.instants.length; i++) {
                        if (d.id === bModule.instants[i].id) {
                            ins = bModule.instants[i];
                        }
                    }

                    var bgColor = "#cfdae1";
                    // console.error("SwitchValue:", ins);
                    if (typeof ins.value === "boolean") {
                        bgColor = ins.value ? (ins.switchFlag <= 1 ? "#3598dc" : "#f7ca18") : "#cfdae1";
                    } else if (typeof ins.value === "number") {
                        bgColor = ins.value > 0 ? (ins.switchFlag <= 1 ? "#3598dc" : "#f7ca18") : "#cfdae1";
                    }

                    var switchLabel = statusModule.append("rect")
                        .attr("x", cx + 4)
                        .attr("y", cy + height / 2 + 4)
                        .attr("width", width - 8)
                        .attr("height", height / 2 - 8)
                        .attr("rx", 6)
                        .attr("ry", 6)
                        .style("fill", bgColor);

                    bModule.switchLabels[d.id] = switchLabel;
                }

                if (d.type === "range") {
                    //StatusColor Drawing
                    var ins;
                    for (var i = 0; i < bModule.instants.length; i++) {
                        if (d.id === bModule.instants[i].id) {
                            ins = bModule.instants[i];
                        }
                    }

                    var bgColor = "#32c5d2";
                    textColor = "#fff";

                    statusModule.append("rect")
                        .attr("x", cx + 4)
                        .attr("y", cy + height / 2 + 4)
                        .attr("width", width - 8)
                        .attr("height", height / 2 - 8)
                        .attr("rx", 6)
                        .attr("ry", 6)
                        .style("fill", bgColor);
                }

                //Bar Drawing
                var barSize = fontSize;
                if (d.name.length > 7) {
                    barSize -= 2 * (d.name.length - 7);
                }
                statusModule.append("text")
                    .attr("x", cx + width / 2)
                    .attr("y", cy + fontSize / 2 + 3)
                    .attr("dy", fontSize / 2)
                    .attr("text-anchor", "middle")
                    .text(d.name)
                    .style("font-size", barSize + "px")
                    //.style("font-weight", "bold")
                    .style("fill", "#fff")
                    .style("stroke-width", "0px");

                //Text Drawing
                var textSize = fontSize - 2;
                if (text.length > 7) {
                    textSize -= 2 * (text.length - 7);
                }
                var valueLabel = statusModule.append("text")
                    .attr("x", cx + width / 2)
                    .attr("y", cy + height / 2 + fontSize / 2 + 2)
                    .attr("dy", fontSize / 2)
                    .attr("text-anchor", "middle")
                    .style("font-size", textSize + "px")
                    //.style("font-weight", "bold")
                    .style("fill", textColor)
                    .style("stroke-width", "0px");

                if (d.type !== "switch" && d.type !== "status") {
                    bModule.valueLabels[d.id] = valueLabel;
                } else if (d.type === "status") {
                    bModule.statusLabels[d.id].value = valueLabel;
                }
                valueLabel.text(text);
            }
        }
    };

    bModule.renderAnimations = function () {
        if (!bModule.isBoilerBurning || !bModule.isOnline) {
            return;
        }

        if (!bModule.svg) {
            return;
        }

        switch (bModule.moduleId) {
            case BOILE_MODULE_OIL:
                renderGasFire("#fire_container");
                renderOilDashes("#dash_container");
                renderGasSmokeDashes("#dash_smoke_container");
                break;
            case BOILE_MODULE_GAS:
                renderGasFire("#fire_container");
                renderGasDashes("#dash_container");
                renderGasSmokeDashes("#dash_smoke_container");
                break;
            case BOILE_MODULE_WATER:
                renderWaterFire("#fire_container");
                renderWaterDashes("#dash_container");
                break;
            case BOILE_MODULE_HEAT_WATER_SYSTEM:
                renderElectricDashes("#dash_container");
                break;
            case BOILE_MODULE_LV:
                renderLvFire("#fire_container");
                renderLvDashes("#dash_container");
                renderLvSmokeDashes("#dash_smoke_container");
                renderLvFan("#fan");
                break;
            case BOILE_MODULE_IRON:
                renderIronFire("#fire_container");
                renderIronDashes("#dash_container");
                renderIronFan("#fan");
                break;
            default:
                renderCoalDashes("#dash_container");
                break;
        }
    };

    var renderCoalDashes = function (id) {
        var dashModule = bModule.svg.select(id);

        var size = 8;
        var sec = 4096;

        var color = "#fff";

        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        // var pathSteam = [
        //     {
        //         cx: 443,
        //         cy: 320
        //     },
        //     {
        //         cx: 443,
        //         cy: 58,
        //         duration: sec / 2
        //     },
        //     {
        //         cx: 690,
        //         cy: 58,
        //         duration: sec / 2
        //     }
        // ];
        // var steam = bModule.dash
        //     .append("circle").attr("cx", pathSteam[0].cx).attr("cy", pathSteam[0].cy).attr("r", size / 2).style("fill", color);
        // d3.selectAll("circle").transition().duration(pathSteam[1].duration).ease(d3.easeLinear).attr("cx", pathSteam[1].cx).attr("cy", pathSteam[1].cy);
        // var dd = function (path) {
        //     d3.select(twizzleLock).transition().duration(pathSteam[1].duration).ease(d3.easeLinear).attr("cx", pathSteam[1].cx).attr("cy", pathSteam[1].cy)
        // };
        // steamDash.remove();

        var dashSteam = function () {
            //d3.active(this).enter()
            dashModule
                .append("circle").attr("cx", 443).attr("cy", 320).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 58)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 690)
                .remove();
        };

        var dashWater = function () {
            dashModule.append("circle").attr("cx", 1200).attr("cy", 135).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 6).ease(d3.easeLinear).attr("cx", 1108)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 440)
                .transition().duration(sec / 16).ease(d3.easeLinear).attr("cx", 1090)
                .remove();

            dashModule.append("circle").attr("cx", 1010).attr("cy", 462).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 5).ease(d3.easeLinear).attr("cx", 924)
                .transition().duration(sec / 5).ease(d3.easeLinear).attr("cy", 570)
                .remove();

            dashModule
                .append("circle").attr("cx", 804).attr("cy", 570).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 238)
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 606)
                .transition().duration(sec / 6).ease(d3.easeLinear).attr("cy", 320)
                .remove();
        };

        var dashSmoke = function () {
            //d3.active(this).enter()
            dashModule
                .append("circle").attr("cx", 630).attr("cy", 402).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 758)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 620)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 1000)
                .remove();
        };

        var dashAir = function () {
            //d3.active(this).enter()
            dashModule
                .append("circle").attr("cx", 270).attr("cy", 716).attr("r", size / 2).style("fill", "#666")
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 436)
                .transition().duration(sec / 4).ease(d3.easeLinear).attr("cy", 650)
                .remove();
        };

        dashModule
            .transition().on("start", function repeat() {
                if (!bModule.isBoilerBurning) {
                    return;
                }
                dashModule
                    .transition().delay(260).on("start", function () {
                        dashSteam();
                        dashWater();
                        dashSmoke();
                        dashAir();
                        repeat();
            });
        });
    };

    var renderOilDashes = function (id) {
        var size = 8;
        var sec = 4096;

        var color = "#fff";

        var dashModule = bModule.svg.select(id);
        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashSteam = function () {
            dashModule
                .append("circle").attr("cx", 611).attr("cy", 340).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 88)
                .transition().duration(sec / 6).ease(d3.easeLinear).attr("cx", 540)
                .remove();
        };

        var dashWater = function () {
            dashModule
                .append("circle").attr("cx", 1200).attr("cy", 377).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 5).ease(d3.easeLinear).attr("cx", 1104)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 626)
                .transition().duration(sec / 16).ease(d3.easeLinear).attr("cx", 1086)
                .remove();

            dashModule
                .append("circle").attr("cx", 1002).attr("cy", 648).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 16).ease(d3.easeLinear).attr("cx", 981)
                .transition().duration(sec / 1).ease(d3.easeLinear).attr("cy", 110)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 672)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 340)
                .remove();


        };

        var dashFuel = function () {
            dashModule
                .append("circle").attr("cx", 130).attr("cy", 514).attr("r", size / 2).style("fill", "#eee")
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 436)
                .remove();
        };

        dashModule
            .transition().on("start", function repeat() {
            dashModule
                .transition().delay(260).on("start", function () {
                dashSteam();
                dashWater();
                dashFuel();
                repeat();
            });
        });
    };

    var renderGasDashes = function (id) {
        var size = 8;
        var sec = 4096;

        var color = "#fff";

        var dashModule = bModule.svg.select(id);
        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashSteam = function () {
            dashModule
                .append("circle").attr("cx", 611).attr("cy", 340).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 88)
                .transition().duration(sec / 6).ease(d3.easeLinear).attr("cx", 540)
                .remove();
        };

        var dashWater = function () {
            dashModule
                .append("circle").attr("cx", 1200).attr("cy", 377).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 5).ease(d3.easeLinear).attr("cx", 1104)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 626)
                .transition().duration(sec / 16).ease(d3.easeLinear).attr("cx", 1086)
                .remove();

            dashModule
                .append("circle").attr("cx", 1002).attr("cy", 648).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 16).ease(d3.easeLinear).attr("cx", 981)
                .transition().duration(sec / 1.5).ease(d3.easeLinear).attr("cy", 230)
                .transition().duration(sec / 8).ease(d3.easeLinear).attr("cx", 922)
                .remove();

            dashModule
                .append("circle").attr("cx", 922).attr("cy", 144).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 8).ease(d3.easeLinear).attr("cx", 981)
                .transition().duration(sec / 12).ease(d3.easeLinear).attr("cy", 110)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 672)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 340)
                .remove();
        };

        var dashFuel = function () {
            dashModule
                .append("circle").attr("cx", 130).attr("cy", 514).attr("r", size / 2).style("fill", "#eee")
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 436)
                .remove();
        };

        dashModule
            .transition().on("start", function repeat() {
            if (!bModule.isBoilerBurning) {
                return;
            }
            dashModule
                .transition().delay(260).on("start", function () {
                dashSteam();
                dashWater();
                dashFuel();
                repeat();
            });
        });
    };

    var renderGasSmokeDashes = function (id) {
        var size = 8;
        var sec = 4096;

        var dashSmokeModule = bModule.svg.select(id);
        if (!dashSmokeModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashSmoke = function () {
            dashSmokeModule
                .append("circle").attr("cx", 896).attr("cy", 330).attr("r", size / 2).style("fill", "#999")
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 80)
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 1064)
                .remove();
        };

        dashSmokeModule
            .transition().on("start", function repeat() {
            if (!bModule.isBoilerBurning) {
                return;
            }
            dashSmokeModule
                .transition().delay(260).on("start", function () {
                dashSmoke();
                repeat();
            });
        });
    };

    var renderWaterDashes = function (id) {
        var waterSize = 6;
        var size = 8;
        var sec = 4096;

        var color = "#fff";

        var dashModule = bModule.svg.select(id);
        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashWaterIn = function () {
            //d3.active(this).enter()
            dashModule
                .append("circle").attr("cx", 400).attr("cy", 220).attr("r", waterSize / 2).style("fill", color)
                .transition().duration(sec / 1.2).ease(d3.easeLinear).attr("cx", 635)
                .transition().duration(sec / 1.8).ease(d3.easeLinear).attr("cy", 330)
                .remove();

            dashModule
                .append("circle").attr("cx", 400).attr("cy", 262).attr("r", waterSize / 2).style("fill", color)
                .transition().duration(sec / 1.2).ease(d3.easeLinear).attr("cx", 615)
                .transition().duration(sec / 2.6).ease(d3.easeLinear).attr("cy", 330)
                .remove();
        };

        var dashWaterOut = function () {
            dashModule.append("circle").attr("cx", 898).attr("cy", 330).attr("r", waterSize / 2).style("fill", color)
                .transition().duration(sec / 1.8).ease(d3.easeLinear).attr("cy", 222)
                .transition().duration(sec / 1.2).ease(d3.easeLinear).attr("cx", 1135)
                .remove();

            dashModule.append("circle").attr("cx", 918).attr("cy", 330).attr("r", waterSize / 2).style("fill", color)
                .transition().duration(sec / 2.6).ease(d3.easeLinear).attr("cy", 260)
                .transition().duration(sec / 1.2).ease(d3.easeLinear).attr("cx", 1135)
                .remove();
        };

        var dashSmoke = function () {
            //d3.active(this).enter()
            dashModule
                .append("circle").attr("cx", 955).attr("cy", 540).attr("r", size / 2).style("fill", "#999")
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 1100)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 350)
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 1200)
                .remove();
        };

        var dashFuel = function () {
            dashModule
                .append("circle").attr("cx", 195).attr("cy", 548).attr("r", size / 2).style("fill", "#eee")
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 420)
                .remove();
        };

        dashModule
            .transition().on("start", function repeat() {
            if (!bModule.isBoilerBurning) {
                return;
            }
            dashModule
                .transition().delay(360).on("start", function () {
                dashWaterIn();
                dashWaterOut();
                dashSmoke();
                dashFuel();
                repeat();
            });
        });
    };

    var renderElectricDashes = function (id) {
        var size = 6;
        var sec = 4096;
        var color = "#fff";

        var dashModule = bModule.svg.select(id);
        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashWater = function () {
            dashModule.append("circle").attr("cx", 263).attr("cy", 270).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 6).ease(d3.easeLinear).attr("cx", 294)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 441)
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 400)
                .remove();

            dashModule.append("circle").attr("cx", 402).attr("cy", 540).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 260)
                .remove();
        };

        dashModule
            .transition().on("start", function repeat() {
            if (!bModule.isBoilerBurning) {
                return;
            }
            dashModule
                .transition().delay(260).on("start", function () {
                dashWater();
                repeat();
            });
        });
    };

    var renderGasFire = function (id) {
        console.info("renderGasFire");
        var fireG = bModule.svg.select(id);
        if (!fireG) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var svgName = "/img/module/boiler_gas_fire.svg";

        var baseX = 540;
        var baseY = 475;

        var fire = fireG.append("svg:image")
            .attr("xlink:href", svgName)
            .attr("width", 120)
            .attr("height", 60)
            .attr("x", baseX)
            .attr("y", baseY);

        var sec = 600;

        var burn = function () {
            fire.transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 180).attr("height", 90).attr("y", baseY - 15)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 120).attr("height", 60).attr("y", baseY);
        };

        d3.interval(function () {
            if (!bModule.isBoilerBurning) {
                fire.remove();
                return;
            }
            burn();
        }, sec);
    };

    var renderWaterFire = function (id) {
        console.info("renderWaterFire");
        var fireG = bModule.svg.select(id);
        if (!fireG) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var svgName = "/img/module/boiler_water_fire.svg";

        var baseX = 660;
        var baseY = 519;

        var fire = fireG.append("svg:image")
            .attr("xlink:href", svgName)
            .attr("width", 100)
            .attr("height", 50)
            .attr("x", baseX)
            .attr("y", baseY);

        var sec = 600;

        var burn = function () {
            fire.transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 160).attr("height", 80).attr("y", baseY - 15)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 100).attr("height", 50).attr("y", baseY);
        };

        d3.interval(function () {
            if (!bModule.isBoilerBurning) {
                fire.remove();
                return;
            }
            burn();
        }, sec);
    };

    //铸铝
    var renderLvDashes = function (id) {
        var waterSize = 6;
        var size = 8;
        var sec = 4096;

        var color = "#fff";

        var dashModule = bModule.svg.select(id);
        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashWaterIn = function () {
            //d3.active(this).enter()

            dashModule
                .append("circle").attr("cx", 455).attr("cy", 138).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 1).ease(d3.easeLinear).attr("cx", 880)
                .remove();

            dashModule
                .append("circle").attr("cx", 880).attr("cy", 420).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 1).ease(d3.easeLinear).attr("cx", 455)
                .remove();


        };

        var dashWaterOut = function () {
            dashModule.append("circle").attr("cx", 895).attr("cy", 500).attr("r", waterSize / 3).style("fill", color)
                .transition().duration(sec / 2.8).ease(d3.easeLinear).attr("cx", 932).attr("cy", 525)
                .transition().duration(sec / 2.8).ease(d3.easeLinear).attr("cx", 980)
                .remove();

        };


        dashModule
            .transition().on("start", function repeat() {
            dashModule
                .transition().delay(360).on("start", function () {
                dashWaterIn();
                dashWaterOut();
                repeat();
            });
        });
    };

    var renderLvSmokeDashes = function (id) {
        var size = 8;
        var sec = 4096;

        var dashSmokeModule = bModule.svg.select(id);
        if (!dashSmokeModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashSmoke = function () {
            dashSmokeModule
                .append("circle").attr("cx", 440).attr("cy", 208).attr("r", size / 2).style("fill", "#999")
                .transition().duration(sec / 1.0).ease(d3.easeLinear).attr("cx", 766)
                .transition().duration(sec / 1.5).ease(d3.easeLinear).attr("cy", 478)
                .transition().duration(sec / 2.5).ease(d3.easeLinear).attr("cx", 955)
                .remove();

            dashSmokeModule
                .append("circle").attr("cx", 440).attr("cy", 208).attr("r", size / 2).style("fill", "#999")
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 533)
                .transition().duration(sec / 1.5).ease(d3.easeLinear).attr("cy", 478)
                .transition().duration(sec / 1).ease(d3.easeLinear).attr("cx", 955)
                .remove();

        };

        dashSmokeModule
            .transition().on("start", function repeat() {
            dashSmokeModule
                .transition().delay(260).on("start", function () {
                dashSmoke();
                repeat();
            });
        });
    };

    var renderLvFan = function (id) {
        var fan = bModule.svg.select(id);
        var fan_inner = fan.select("#fan_inner");

        var sec = 600;

        var fanmove = function () {
            fan_inner.remove();
            fan.append("svg:image").attr("xlink:href", "../img/fan.gif")
                .attr("x", 320)
                .attr("y", 200);

        };

        fanmove();

    };


    var renderLvFire = function (id) {
        var fireG = bModule.svg.select(id);
        if (!fireG) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var svgName = "../img/module/boiler_gas_fire2.svg";

        var baseX = 500;
        var baseY = 190;

        var fire = fireG.append("svg:image")
            .attr("xlink:href", svgName)
            .attr("width", 180)
            .attr("height", 90)
            .attr("x", baseX)
            .attr("y", baseY);

        var sec = 600;

        var burn = function () {
            fire.transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 270).attr("height", 135).attr("y", baseY - 15)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 180).attr("height", 90).attr("y", baseY);
        };

        d3.interval(function () {
            burn();
        }, sec);

    };


    //铸铁
    var renderIronDashes = function (id) {
        var waterSize = 6;
        var size = 8;
        var sec = 4096;

        var color = "#fff";

        var dashModule = bModule.svg.select(id);
        if (!dashModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashWaterIn = function () {
            //d3.active(this).enter()

            dashModule
                .append("circle").attr("cx", 890).attr("cy", 290).attr("r", size / 1.2).style("fill", "#ed6b75")
                .transition().duration(sec / 0.8).ease(d3.easeLinear).attr("cx", 112)
                .remove();

            dashModule
                .append("circle").attr("cx", 268).attr("cy", 610).attr("r", size / 2).style("fill", color)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 390)
                .remove();

            dashModule
                .append("circle").attr("cx", 938).attr("cy", 350).attr("r", size / 1.2).style("fill", "#ed6b75")
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cy", 60)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("cx", 450)
                .remove();
        };

        var dashWaterOut = function () {
            dashModule.append("circle").attr("cx", 895).attr("cy", 500).attr("r", waterSize / 3).style("fill", color)
                .transition().duration(sec / 2.8).ease(d3.easeLinear).attr("cx", 932).attr("cy", 525)
                .transition().duration(sec / 2.8).ease(d3.easeLinear).attr("cx", 980)
                .remove();

        };


        dashModule
            .transition().on("start", function repeat() {
            dashModule
                .transition().delay(360).on("start", function () {
                dashWaterIn();
//              dashWaterOut();
                repeat();
            });
        });
    };

    var renderIronSmokeDashes = function (id) {
        var size = 8;
        var sec = 4096;

        var dashSmokeModule = bModule.svg.select(id);
        if (!dashSmokeModule) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var dashSmoke = function () {
            dashSmokeModule
                .append("circle").attr("cx", 440).attr("cy", 208).attr("r", size / 2).style("fill", "#999")
                .transition().duration(sec / 1.0).ease(d3.easeLinear).attr("cx", 766)
                .transition().duration(sec / 1.5).ease(d3.easeLinear).attr("cy", 478)
                .transition().duration(sec / 2.5).ease(d3.easeLinear).attr("cx", 955)
                .remove();

            dashSmokeModule
                .append("circle").attr("cx", 440).attr("cy", 208).attr("r", size / 2).style("fill", "#999")
                .transition().duration(sec / 3).ease(d3.easeLinear).attr("cx", 533)
                .transition().duration(sec / 1.5).ease(d3.easeLinear).attr("cy", 478)
                .transition().duration(sec / 1).ease(d3.easeLinear).attr("cx", 955)
                .remove();

        };

        dashSmokeModule
            .transition().on("start", function repeat() {
            dashSmokeModule
                .transition().delay(260).on("start", function () {
                dashSmoke();
                repeat();
            });
        });
    };

    var renderIronFan = function (id) {
        var fan = bModule.svg.select(id);
        var fan_inner = fan.select("#fan_inner");

        var sec = 600;

        var fanmove = function () {
            fan_inner.remove();
            fan.append("svg:image").attr("xlink:href", "../img/fan.gif")
                .attr("x", 1040)
                .attr("y", 435);

        };

        fanmove();

    };



    var renderIronFire = function (id) {
        var fireG = bModule.svg.select(id);
        if (!fireG) {
            console.warn("There IS NO " + id + "!");
            return;
        }

        var svgName = "../img/module/boiler_water_fire_2.svg";

        var baseX = 660;
        var baseY = 380;

        var fire = fireG.append("svg:image")
            .attr("xlink:href", svgName)
            .attr("width", 240)
            .attr("height", 120)
            .attr("x", baseX)
            .attr("y", baseY);

        var sec = 600;

        var burn = function () {
            fire.transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 380).attr("height", 180).attr("y", baseY - 15).attr("x", baseX - 140)
                .transition().duration(sec / 2).ease(d3.easeLinear).attr("width", 240).attr("height", 120).attr("y", baseY).attr("x", baseX);
        };

        d3.interval(function () {
            burn();
        }, sec);
    };

    bModule.updateStatusLabels = function () {
        if (!bModule.boiler) {
            if (!$rootScope.boiler) {
                return
            }

            bModule.boiler = $rootScope.boiler;
        }

        var isTerminalConnected = (bModule.boiler.Terminal && bModule.boiler.Terminal.IsOnline) || bModule.isBoilerBurning;

        for (var i in bModule.statusLabels) {
            var statusLabel = bModule.statusLabels[i];
            var text = "";
            var bgColor = "#32c5d2";

            switch (parseInt(i, 10)) {
                case 1:
                    text = isTerminalConnected ? "已连接" : "未连接";
                    bgColor = isTerminalConnected ? "#32c5d2" : "#bfcad1";
                    break;
                case 2:
                    text = bModule.isBoilerBurning ? "已点燃" : "未点燃";
                    bgColor = bModule.isBoilerBurning ? "#32c5d2" : "#e7505a";
                    break;
                case 3:
                    switch (bModule.boiler.alarmLevel) {
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

    bModule.updateLabels = function () {
        // $log.error("updateWaterText()", bModule.valueLabels, new Date());
        for (var i in bModule.instants) {
            var ins = bModule.instants[i];
            if (bModule.valueLabels[ins.id] &&
                ins.category !== 11) {
                bModule.valueLabels[ins.id].text(ins.value + ins.unit);
            }

            if (bModule.switchLabels[ins.id]) {
                var bgColor = "#cfdae1";
                // console.error("SwitchValue:", ins);
                if (typeof ins.value === "boolean") {
                    bgColor = ins.value ? (ins.switchFlag <= 1 ? "#3598dc" : "#f7ca18") : "#cfdae1";
                } else if (typeof ins.value === "number") {
                    bgColor = ins.value > 0 ? (ins.switchFlag <= 1 ? "#3598dc" : "#f7ca18") : "#cfdae1";
                }

                bModule.switchLabels[ins.id].style("fill", bgColor);
            }
        }
    };
});

var hasInitModule = false;
var bModule;

const BOILE_MODULE_COAL_DOUBLE = 1;
const BOILE_MODULE_OIL = 2;
const BOILE_MODULE_GAS = 3;
const BOILE_MODULE_WATER = 4;
const BOILE_MODULE_COAL_WATER = 6;
const BOILE_MODULE_HEAT_WATER_SYSTEM = 11;
const BOILE_MODULE_LV = 12;
const BOILE_MODULE_IRON = 13;
const BOILE_MODULE_HEAT_STEAM_SYSTEM = 14;
