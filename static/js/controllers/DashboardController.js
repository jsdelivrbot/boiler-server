angular.module('BoilerAdmin').controller('DashboardController', function($rootScope, $scope, $http, $filter, $timeout, $log, $uibModal, $document, $location, $state, $anchorScroll, moment, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    bMonitor = this;
    bMonitor.isDone = false;

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    bMonitor.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
    //.withOption('rowCallback', rowBoilerCallback);

    bMonitor.mapOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('simple')
        .withOption('searching', false)
        .withOption('rowCallback', bMonitor.mapRowClicked);

    bMonitor.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5).notSortable()
    ];

    bMonitor.mapColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0).notSortable()
    ];

    bMonitor.mapRowClicked = function (boiler) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        console.warn("Click Row:", boiler);
        if (!boiler.Address) {
            console.warn("Boiler Has No Address!");
            return;
        }

        console.warn("Ready to Move:", boiler.Address.Longitude, boiler.Address.Latitude);
        var point = new BMap.Point(boiler.Address.Longitude, boiler.Address.Latitude);
        bMonitor.bMap.setZoom(14);
        bMonitor.bMap.panTo(point);
    };

    bMonitor.getBoilers = function () {
        bMonitor.datasource = $rootScope.boilers;
        if (!$rootScope.boilers) {
            return;
        }
        for (var i = 0; i < bMonitor.datasource.length; i++) {
            var d = bMonitor.datasource[i];
            d.num = i;
            d.name = d.Name;
            if (d.name.length > 12) {
                d.name = d.name.substring(0, 10) + '...';
            }
            // bMonitor.fetchStatus(d);
            if (d.Status) {
                d.isBurning = d.Status.IsBurning > 0;
            }
        }

        bMonitor.initSearch();
        bMonitor.filterBoilers();
        setTimeout(function () {
            App.stopPageLoading();
        }, 800);

        bMonitor.isDone = $rootScope.boilers.length > 0;
    };

    $rootScope.$watch('boilers', function () {
        // $log.warn("$rootScope.$watch.boilers: ", $rootScope.boilers);
        if (!$rootScope.boilers || typeof $rootScope.boilers !== 'object') {
            return;
        }
        bMonitor.getBoilers();
        bMonitor.initAmChartPie();
    });

    $rootScope.$watch('locations', function () {
        bMonitor.initSearch();
    });

    bMonitor.getRuntimeCount = function () {
        $http.get('/boiler_runtime_count/')
            .then(function (res) {
                console.warn("boiler_runtime_count resp:", res);
                bMonitor.runtimeCount = res.data;
            });
    };

    bMonitor.fetchTotal = function () {
        var week = 4;

        $http.get('/boiler_status_running/')
            .then(function (res) {
                console.info("boiler_status_running resp:", res);
                var total = 0;
                for (var i = 0; i < Object.values(res.data).length; i++) {
                    var due = Object.values(res.data)[i];
                    total += due;
                }
                console.info("total duration: ", due);
                var duraText = '';

                var duration = total / 1000 / 1000 / 1000;   //sec
                var dHour = Math.floor(duration / 60 / 60);

                duraText += dHour > 0 ? dHour + ' 时' : '';
                /*
                var dHour = Math.floor(duration / 60 / 60);
                duration -= dHour * 60 * 60;
                var dMin = Math.floor(duration / 60);
                duration -= dMin * 60;

                duraText += dHour > 0 ? dHour + '小时' : '';
                duraText += dMin + '分';
                duraText += duration.toFixed(2) + '秒';
                */

                bMonitor.runningTotal = duraText;
        });

        $http.get('/boiler_evaporate_rank/').then(function (res) {
            console.warn("boiler_evaporate_rank RESP:", res);
            var chart = new AmCharts.AmSerialChart();
            chart.theme = AmCharts.themes.light;
            chart.color = "#333";
            chart.language = "zh";
            chart.dataDateFormat = 'MMM DDD';
            //chart.dataProvider = dataProvider,
            chart.valueAxes = [{
                stackType: "none",
                position: "left",
                maximum: 105,
                showLastLabel: false,
                unit: '%',
                title: '达标率'
            }];
            chart.colors= [
                "#67b7dc",
                "#c4e479",
                "#84b761",
                "#cc4748",
                "#cd82ad",
                "#2f4074",
                "#448e4d",
                "#b7b83f",
                "#b9783f",
                "#b93e3d",
                "#913167"
            ];
            chart.legend = {
                horizontalGap: 10,
                useGraphSettings: true,
                markerSize: 10,
                valueWidth: 50,
            };
            // chart.legend = {
            //     position: "absolute",
            //     top: "10px",
            //     right: "10px",
            //     valueText: "[[value]]",
            //     valueWidth: 100,
            //     valueAlign: "left",
            //     equalWidths: false,
            //     periodValueText: "total: [[value.sum]]"
            //
            // };
            chart.startDuration = 1;
            // chart.graphs: graphs,
            chart.plotAreaFillAlphas = 0.1;
            chart.depth3D = 26;
            chart.angle = 45;
            chart.categoryField = "evaporate";
            chart.categoryAxis = {
                // minPeriod: "mm",
                // parseDates: true,
                // equalSpacing: true,
                axisAlpha: 0.2,
                gridPosition: "start",
                title: "燃煤锅炉、生物质锅炉　　　　　　　　　　燃油锅炉、燃气锅炉"
            };
            // chart.allLabels = [
            //     {
            //         text: "燃油锅炉、燃气锅炉",
            //         align: "right",
            //         size: 12,
            //         bold: true,
            //         x: '92%',
            //         y: 475
            //     }
            // ];
            chart.export = {
                enabled: true
            };

            var ids = ['c0', 'c1', 'c2', 'c3', 'c4', 'g0', 'g1'];
            var evaporates = ['D≤1<br>(≥61%)', '1＜D≤2<br>(≥69%)', '2＜D≤8<br>(≥71%)', '8＜D≤20<br>(≥72%)', 'D＞20<br>(≥72%)', 'D≤2<br>(≥79%)', 'D＞2<br>(≥81%)'];

            var dataProvider = [];

            for (var i = 0; i < ids.length; i++) {
                var data = {};
                data.id = ids[i];
                data.evaporate = evaporates[i];
                data.countSuccess = 0;
                data.countFailed = 0;
                data.percentSuccess = 0.0;
                data.percentFailed = 0.0;
                data.colorSuccess = "#67b7dc";
                data.colorFailed = "#fdd400";

                var items = $filter('filter')(res.data, function (item) {
                    if (item['evaporate_id'] === data.id) {
                        return true;
                    }
                    return false;
                });

                if (items && items.length > 0) {
                    for (var j = 0; j < items.length; j++) {
                        var it = items[j];
                        if (it.rank === 'success') {
                            data.countSuccess = parseInt(it.count);
                        } else {
                            data.countFailed = parseInt(it.count);
                        }
                    }
                }

                if (data.countSuccess + data.countFailed > 0) {
                    data.percentSuccess = (data.countSuccess * 100) / (data.countSuccess + data.countFailed);
                    data.percentFailed = 100 - data.percentSuccess;

                    data.percentSuccess = data.percentSuccess.toFixed(1);
                    data.percentFailed = data.percentFailed.toFixed(1);
                }

                dataProvider.push(data);

                //colors: ["#67b7dc", "#fdd400", "#84b761", "#cc4748", "#cd82ad", "#2f4074", "#448e4d", "#b7b83f", "#b9783f", "#b93e3d", "#913167"]
            }

            var status = ['Success', 'Failed'];

            for (var is = 0; is < status.length; is++) {
                var st = status[is];

                var graph = new AmCharts.AmGraph();
                graph.fillAlphas = 0.66;
                graph.lineAlpha = 0.2;
                graph.title = st;
                graph.fillColorsField= "color";
                //graph.title = "[[date + num]]";
                // graph.highField = "count" + st;
                graph.labelText = "[[ high ]]";
                graph.labelFunction = function (graphDataItem) {
                    var field = "count" + (graphDataItem.color === "#67b7dc" ? "Success" : "Failed");
                    var text = graphDataItem.dataContext[field];
                    // console.warn("labelText", graphDataItem, text, field, st);
                    return text + " 台";
                };
                graph.type = "column";
                // graph.colorField = "color" + st;
                graph.valueField = "percent" + st;
                graph.balloonText = "[[value]]" + " %";
                // graph.balloonFunction = balloneText;

                chart.addGraph(graph);
            }

            chart.dataProvider = dataProvider;
            // console.error("Rank dataProvider:", dataProvider);

            if ("undefined" != typeof AmCharts && 0 !== $("#dashboard_amchart_3d").size()) {
                chart.write("dashboard_amchart_3d");
            }
        });

        /*
        $http.get('/boiler_runtime_daily_total/')
            .then(function (res) {
                console.warn("boiler_runtime_daily_total RESP:", res);
                var chart = new AmCharts.AmSerialChart();
                chart.theme = "light";
                chart.language = "zh";
                chart.dataDateFormat = 'MMM DDD';
                //chart.dataProvider = dataProvider,
                chart.valueAxes = [{
                    stackType: "3d",
                    unit: "t",
                    position: "left"
                }];
                chart.startDuration = 1;
                // chart.graphs: graphs,
                chart.plotAreaFillAlphas = 0.1;
                chart.depth3D = 120;
                chart.angle = 60;
                chart.categoryField = "num";
                chart.categoryAxis = {
                    // minPeriod: "mm",
                    // parseDates: true,
                    // equalSpacing: true,
                    axisAlpha: 0.2,
                    gridPosition: "start"
                    //title: "本月平均每日流量"
                };
                chart.export = {
                    enabled: true
                };

                var dataProvider = [];

                var weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六'];

                var balloneText = function (graphDataItem, graph) {
                    //console.error("balloneText", graphDataItem, graph);
                    var dateText = graphDataItem.dataContext["date" + graph.index];
                    var value = graphDataItem.values.value;
                    return dateText + "<br><b>日均累计流量: " + value + "t</b>";
                };

                for (var num = 0; num < week; num++) {
                    for (var i = 0; i < 7; i++) {
                        var data = {};
                        if (i < dataProvider.length) {
                            data = dataProvider[i];
                        }

                        var aDay = new Date();
                        aDay.setHours(0);
                        aDay.setMinutes(0);
                        aDay.setSeconds(0);
                        aDay.setDate(aDay.getDate() - (week - num) * 7 + i);
                        var mDate = moment(aDay);

                        data.num = weekdays[aDay.getDay()];
                        data['date' + num] = mDate.format("MMM DD, YYYY");
                        data['flow' + num] = "0.00";
                        data['color' + num] = "#999999";

                        for (var j = 0; j < res.data.length; j++) {
                            var d = res.data[j];
                            var date = new Date(d.Date);

                            if (date.getDate() == aDay.getDate()) {
                                data['flow' + num] = d.Flow.toFixed(2);
                                if (d.Flow >= 44) {
                                    data['color' + num] = "#67b7dc";
                                } else {
                                    data['color' + num] = "#fdd400";
                                }
                                break;
                            }
                        }

                        if (i >= dataProvider.length) {
                            dataProvider.push(data);
                        }
                    }

                    var graph = new AmCharts.AmGraph();
                    graph.fillAlphas = 0.66;
                    graph.lineAlpha = 0.2;
                    //graph.title = "[[date + num]]";
                    graph.type = "column";
                    graph.colorField = "color" + num;
                    graph.valueField = "flow" + num;
                    graph.balloonFunction = balloneText;

                    chart.addGraph(graph);

                    //colors: ["#67b7dc", "#fdd400", "#84b761", "#cc4748", "#cd82ad", "#2f4074", "#448e4d", "#b7b83f", "#b9783f", "#b93e3d", "#913167"]
                }

                chart.dataProvider = dataProvider;

                console.error("dataProvider:", dataProvider);

                if ("undefined" != typeof AmCharts && 0 !== $("#dashboard_amchart_3d").size()) {
                    chart.write("dashboard_amchart_3d");
                }
            });
            */
    };

    /**
     * Search Section
     */
    bMonitor.initSearch = function () {
        bMonitor.provinces = $rootScope.locations;
        bMonitor.citis = [];
        bMonitor.regions = [];

        bMonitor.organizations = [{Uid: '', name: '所属企业（不限）'}];

        bMonitor.aBurning = undefined;

        bMonitor.evaporates = [
            {id: 0, Text: '额定蒸发量（不限）'},
            {id: 1, Text: 'D≤1'},
            {id: 2, Text: '1<D≤2'},
            {id: 3, Text: '2<D≤8'},
            {id: 4, Text: '8<D≤20'},
            {id: 5, Text: 'D>20'}
        ];
        bMonitor.mediums = [{Name: '锅炉介质（不限）'}];
        bMonitor.forms = [{Name: '锅炉型态（不限）'}];
        bMonitor.fuels = [{Name: '锅炉燃料（不限）'}];

        for (var i = 0; $rootScope.organizations && i < $rootScope.organizations.length; i++) {
            var org = $rootScope.organizations[i];
            if (bMonitor.organizations.indexOf(org) > -1) {
                continue;
            }
            bMonitor.organizations.push(org);
        }

        for (var i = 0; $rootScope.boilerMediums && i < $rootScope.boilerMediums.length; i++) {
            var med = $rootScope.boilerMediums[i];
            bMonitor.mediums.push(med);
        }

        for (var i = 0; $rootScope.boilerForms && i < $rootScope.boilerForms.length; i++) {
            var form = $rootScope.boilerForms[i];
            if (form.Id === 0 || bMonitor.forms.indexOf(form) > -1) {
                continue;
            }

            bMonitor.forms.push(form);
        }

        for (var i = 0; $rootScope.fuelTypes && i < $rootScope.fuelTypes.length; i++) {
            var fuel = $rootScope.fuelTypes[i];
            if (fuel.Id === 0 || fuel.Id >= 5 || bMonitor.fuels.indexOf(fuel) > -1) {
                continue;
            }
            bMonitor.fuels.push(fuel);
        }

        var localCount = function (locations) {
            //console.warn("localCount", locations);
            if (!locations) {
                return;
            }
            for (var i = 0; i < locations.length; i++)  {
                var local = locations[i];
                var matchNum = 0;
                $filter('filter')(bMonitor.datasource, function (item) {
                    var locationId = (!item.Address || !item.Address.Location) ? 0 : item.Address.Location.LocationId;
                    if (locationId === local.LocationId ||
                        Math.floor(locationId / 100) === local.LocationId ||
                        Math.floor(locationId / 10000) === local.LocationId) {
                        matchNum++;
                        return true;
                    }
                    return false;
                });
                local.count = matchNum;
                if (local.LocationId !== 0) {
                    local.name = local.Name + ' - ' + local.count;
                } else {
                    local.name = '所在区域';
                }

                if (local.SuperId === 0) {
                    localCount(local.cities);
                }

                if (local.SuperId > 0 && local.SuperId < 100) {
                    localCount(local.regions);
                }
            }
        };

        //localCount(bMonitor.provinces);
        bMonitor.aProvince = null;
        if (bMonitor.province && bMonitor.provinces.length > 0) {
            bMonitor.provinces[0].Name = '所在区域';
            bMonitor.aProvince = bMonitor.provinces[0];

        }
        bMonitor.aCity = null;
        bMonitor.aRegion = null;

        bMonitor.aEvaporate = bMonitor.evaporates[0];
        bMonitor.aForm = bMonitor.forms[0];
        bMonitor.aMedium = bMonitor.mediums[0];
        bMonitor.aOrg = bMonitor.organizations[0];
        //bMonitor.aProvince = bMonitor.provinces[0];
        bMonitor.aFuel = bMonitor.fuels[0];

        bMonitor.aLocation = null;
        bMonitor.aQuery = "";

        console.info("bMonitor.aOrg:", bMonitor.aOrg);
    };

    bMonitor.fetchStatus = function (boiler) {
        $http.get('/boiler/state/is_burning/?boiler=' + boiler.Uid)
            .then(function (res) {
                // console.error("Fetch Status Resp:", res.data, boiler.Name);
                boiler.isBurning = (res.data.value === "true");
            }, function (err) {
                console.error('Fetch Status Err!', err);
            });
        $http.get('/boiler/state/is_Online/?boiler=' + boiler.Uid)
            .then(function (res) {
                // console.error("Fetch Status Resp:", res.data, boiler.Name);
                boiler.isOnline = res.data.IsOnline;
            }, function (err) {
                console.error('Fetch Status Err!', err);
            });

    };

    bMonitor.fetchThumbParam = function (boiler) {
        // console.error("boiler:", boiler, boiler.TerminalCode);
        // var rtmQ = [1201, 1015, 1002, 1202];
        // if ($state.current.name !== 'monitor.thumb') {
        //     // $log.error("$state.current.name !== 'monitor.thumb'");
        //     return;
        // }

        if (!bMonitor.pagedItems[bMonitor.currentPage]) {
            return;
        }
        var isMatched = false;
        for (var i = 0; i < bMonitor.pagedItems[bMonitor.currentPage].length; i++) {
            var b = bMonitor.pagedItems[bMonitor.currentPage][i];
            if (b.Uid === boiler.Uid) {
                isMatched = true;
                break;
            }
        }

        if (!isMatched) {
            return;
        }

        var rtmQ = [];

        switch (boiler.TerminalCode) {
            case '680055':
                rtmQ = [1002, 1015];
                break;
            case '680082':
            case '680085':
            case '680096':
            case '680120':
                rtmQ = [1021, 1080];
                break;
            case '680064':
                rtmQ = [1001, 1015];
                break;
            case '680500':
            case '680053':
            case '680501':
            case '680502':
                rtmQ = [1096, 1098, 1090, 1094];
                break;
        }

        $http.post('/boiler_runtime_instants/?scope=thumb', {
            uid: boiler.Uid,
            runtimeQueue: rtmQ
        }).then(function (res) {
            boiler.imgName = function() {
                var imgName = boiler.Form.Id === 101 ? 'fb' : 'coalsingle';
                if (boiler.Fuel && boiler.Fuel.Type) {
                    switch (boiler.Fuel.Type.Id) {
                        case 1:
                        case 4:
                            if (boiler.Form.Id === 201 || boiler.Form.Id === 203) {
                                imgName = 'coalsingle';
                            } else if (boiler.Form.Id === 202 || boiler.Form.Id === 204) {
                                imgName = 'coaldouble';
                            } else if (boiler.Form.Id === 101) {
                                imgName = 'fb';
                            }
                            break;
                        case 2:
                        case 3:
                            if (boiler.Form.Id === 101) {
                                imgName = 'gasboiler_v';
                            } else {
                                imgName = 'gasboiler';
                            }
                    }
                }

                if (boiler.Form.Id === 205 ||
                    (boiler.Medium.Id === 2 && (boiler.Fuel.Type.Id === 2 || boiler.Fuel.Type.Id === 3))) {
                    imgName = 'boilerwater';
                }

                if (boiler.Form.Id === 1003) {
                    imgName = 'gasboiler_v';
                }

                return imgName;
            };
            //console.log("Res:", res.data);
            // var isBurning = function() {
            //     var totalLen = 0;
            //     for (i = 0; i < res.data.length; i++ ) {
            //         var d = res.data[i];
            //         if (d && d.IsValid) {
            //             totalLen ++;
            //         }
            //     }
            //
            //     return totalLen > 0;
            // };
            //
            // boiler.isBurning = isBurning();
            boiler.alarmLevel = boiler.isBurning ? 0 : -1;
            boiler.img = boiler.imgName() + (boiler.isBurning ? '.gif' : '.png');

            var runtime = [[], []];

            for (var i = 0; i < Math.min(res.data.length, 4); i++) {
                var d = res.data[i];
                var value;
                var name = d.ParameterName;
                var alarmLevel = boiler.isBurning ? 0 : -1;

                switch (d.Parameter) {
                    case 1021:
                        name = "环境温度";
                        break;
                    case 1090:
                        name = "省煤器#1出口";
                        break;
                    case 1091:
                        name = "省煤器#2出口";
                        break;
                    case 1202:
                        name = "过量空气系数";
                        break;
                }

                if (boiler.isBurning) {
                    value = d.Value;
                    alarmLevel = d.AlarmLevel;

                    value += " " + d.Unit;
                } else {
                    value = "-";
                }

                if (alarmLevel > boiler.alarmLevel) {
                    boiler.alarmLevel = alarmLevel;
                }

                runtime[i % 2].push({
                    name: name,
                    value: value,
                    alarmLevel: alarmLevel
                });
            }

            boiler.runtime = runtime;

        }, function (err) {
            // alert('Fetch Err!' + err.status + " | " + err.data);

        }).then(function () {
            // console.error("FINALLY bMonitor.fetchThumbParam(boiler):", boiler);

            setTimeout(function () {
                bMonitor.fetchStatus(boiler);
                bMonitor.fetchThumbParam(boiler);
            }, 15000);
        });
    };

    bMonitor.changeProvince = function () {
        bMonitor.aLocation = bMonitor.aProvince;
        bMonitor.filterBoilers();
    };

    bMonitor.changeCity = function () {
        bMonitor.aLocation = bMonitor.aCity;
        bMonitor.filterBoilers();
    };

    bMonitor.changeRegion = function () {
        bMonitor.aLocation = bMonitor.aRegion;
        bMonitor.filterBoilers();
    };

    bMonitor.changeBurning = function () {
        console.warn("Burning:", bMonitor.isBurning);
    };

    bMonitor.filterBoilers = function () {
        var items = bMonitor.filterLocation(bMonitor.datasource);
        // console.warn("items_location:", items.length);
        items = bMonitor.filterOrganization(items);
        // console.warn("items_organization:", items.length);
        items = bMonitor.filterBurning(items);
        // console.warn("items_burning:", items.length);
        items = bMonitor.filterForm(items);
        // console.warn("items_medium:", items.length);
        items = bMonitor.filterFuel(items);
        // console.warn("items_form:", items.length);
        items = bMonitor.filterEvaporate(items);
        // console.warn("items_form:", items.length);
        items = bMonitor.filterQuery(items);
        // console.warn("items_query:", items.length);

        bMonitor.filterLen = items.length;
        // take care of the sorting order
        if (sortingOrder !== '') {
            items = $filter('orderBy')(items, sortingOrder, reverse);
        }

        bMonitor.filteredItems = items;

        // bMonitor.currentPage = 1;
        // bMonitor.pageSize = itemsPerPage;
        // bMonitor.maxSize = pageRange;
        // bMonitor.totalItems = bMonitor.filteredItems.length;

        // now group by pages
        bMonitor.groupToPages();
        // and init BMap
        bMonitor.initBap();
    };

    var sortingOrder = '';
    var reverse = false;
    bMonitor.filteredItems = [];
    var groupedItems = [];
    var itemsPerPage = 4;
    var pageRange = 10;

    bMonitor.pagedItems = [];
    bMonitor.currentPage = 0;
    bMonitor.filterLen = 0;
    //bMonitor.rangedPages = [];

    bMonitor.matchNum = 0;

    var searchMatch = function (haystack, needle) {
        bMonitor.needleNull = "NOT NULL";
        bMonitor.haystack = haystack + ":" + typeof(haystack);
        bMonitor.needle = needle;
        var isMatch = false;
        if (!needle) {
            bMonitor.needleNull = "NULL!!";
            return true;
        }
        if (!haystack) {
            return false;
        }

        switch (typeof haystack) {
            case 'string':
                isMatch = haystack.toLowerCase().indexOf(needle.toLowerCase()) !== -1;
                break;
            case 'object':
                var keys = Object.keys(haystack);
                for (var i = 0; i < keys.length; i++) {
                    if (searchMatch(haystack[keys[i]], needle)) {
                        isMatch = true;
                        break;
                    }
                }
                break;
            default:
                break;
        }

        //isMatch = haystack.toLowerCase().indexOf(needle.toLowerCase()) !== -1;
        //alert(haystack + ", " + needle + ": " + isMatch);
        return isMatch;
    };

    // init the filtered items
    bMonitor.filterQuery = function (boilers) {
        //alert('ready to search: ' + bMonitor.query);

        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            //alert('Item: ' + Object.keys(item));
            if (searchMatch(item, bMonitor.aQuery)) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;
        //alert('runtime.filteredItems Searched: ' + filteredItems.length);
        //alert('runtime.filteredItems Searched: ' + runtime.filteredItems.length + " | " + matchNum);

        return items;
    };

    bMonitor.filterLocation = function (boilers) {
        //alert('ready to search: ' + location.LocationId + ", " + location.Name);
        var matchNum = 0;
        var locationId = !bMonitor.aLocation ? 0 : bMonitor.aLocation.LocationId;
        var items = $filter('filter')(boilers, function (item) {
            //alert('Item: ' + Object.keys(item));
            if (!item.Address ||
                item.Address.Location.LocationId === locationId ||
                Math.floor(item.Address.Location.LocationId / 100) === locationId ||
                Math.floor(item.Address.Location.LocationId / 10000) === locationId ||
                locationId === 0) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;

        return items;
    };

    bMonitor.filterOrganization = function (boilers) {
        //console.warn('ready to filterOrganization: ', bMonitor.aOrg);
        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            //console.warn('Item: ' + Object.keys(item));
            if (!bMonitor.aOrg || !bMonitor.aOrg.Uid || bMonitor.aOrg.Uid.length === 0) {
                matchNum++;
                return true;
            }

            if (item.Factory && item.Factory.Uid === bMonitor.aOrg.Uid ||
                item.Enterprise && item.Enterprise.Uid === bMonitor.aOrg.Uid ||
                item.Installed && item.Installed.Uid === bMonitor.aOrg.Uid) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;

        return items;
    };

    bMonitor.filterBurning = function (boilers) {
        //console.warn('ready to filterOrganization: ', bMonitor.aOrg);
        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            //console.warn('Item: ' + Object.keys(item));
            if (typeof bMonitor.aBurning === 'undefined') {
                matchNum++;
                return true;
            }

            if (bMonitor.aBurning == item.isBurning) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;

        return items;
    };

    bMonitor.filterMedium = function (boilers) {
        //alert('ready to search: ' + location.LocationId + ", " + location.Name);
        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            //console.warn('Item: ', item);
            if ((item.Medium.Id === bMonitor.aMedium.Id) ||
                !bMonitor.aMedium || !bMonitor.aMedium.Id) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;

        return items;
    };

    bMonitor.filterForm = function (boilers) {
        //alert('ready to search: ' + location.LocationId + ", " + location.Name);
        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            //alert('Item: ' + Object.keys(item));
            if ((item.Form.Id === bMonitor.aForm.Id) ||
                !bMonitor.aForm || !bMonitor.aForm.Id) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;

        return items;
    };

    bMonitor.filterFuel = function (boilers) {
        //alert('ready to search: ' + location.LocationId + ", " + location.Name);
        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            //console.warn('filterFuel Item: ' + item);
            if ((item.Fuel.Type.Id === bMonitor.aFuel.Id) ||
                !bMonitor.aFuel || !bMonitor.aFuel.Id || !bMonitor.aFuel.Id === 0) {
                matchNum++;
                return true;
            }
            return false;
        });
        bMonitor.matchNum = matchNum;

        return items;
    };

    bMonitor.filterEvaporate = function (boilers) {
        //console.warn('filterEvaporate: ', bMonitor.aEvaporate);
        var matchNum = 0;
        var items = $filter('filter')(boilers, function (item) {
            // console.warn('filterEvaporate Item: ', item);
            if (!bMonitor.aEvaporate || bMonitor.aEvaporate.id === 0) {
                matchNum++;
                return true;
            }

            var evaporate = item.EvaporatingCapacity;

            switch (bMonitor.aEvaporate.id) {
                case 1:
                    if (evaporate <= 1) {
                        //console.warn("Boiler Evapo &: ", evaporate, "<= 1 Get!");
                        matchNum++;
                        return true;
                    }
                    break;
                case 2:
                    if (evaporate > 1 && evaporate <= 2) {
                        //console.warn("Boiler Evapo &: ", evaporate, "1<D≤2 Get!");
                        matchNum++;
                        return true;
                    }
                    break;
                case 3:
                   // console.warn("Filter Evapo &: ", monitor.data.filter.evaporate_idx);
                    if (evaporate > 2 && evaporate <= 8) {
                        //console.warn("Boiler Evapo &: ", evaporate, "2<D≤8 Get!");
                        matchNum++;
                        return true;
                    }
                    break;
                case 4:
                    if (evaporate > 8 && evaporate <= 20) {
                        //console.warn("Boiler Evapo &: ", evaporate, "8<D≤20 Get!");
                        matchNum++;
                        return true;
                    }
                    break;
                case 5:
                    if (evaporate > 20) {
                        //console.warn("Boiler Evapo &: ", evaporate, "> 20 Get!");
                        matchNum++;
                        return true;
                    }
                    break;
                default:
                    break;
            }

            return false;
        });

        bMonitor.matchNum = matchNum;

        return items;
    };

    // calculate page in place
    bMonitor.groupToPages = function () {
        bMonitor.pagedItems = [];
        for (var i = 0; i < bMonitor.filteredItems.length; i++) {
            if (i % itemsPerPage === 0) {
                bMonitor.pagedItems[Math.floor(i / itemsPerPage)] = [bMonitor.filteredItems[i]];
            } else {
                bMonitor.pagedItems[Math.floor(i / itemsPerPage)].push(bMonitor.filteredItems[i]);
            }
        }
        //alert('runtime.pagedItems after search! ' + runtime.pagedItems.length);
    };

    bMonitor.range = function () {
        var ret = [];
        var length = bMonitor.pagedItems.length;
        var startPage = Math.floor(bMonitor.currentPage / pageRange) * pageRange;
        //alert('Paging: ' + runtime.currentPage + "|" + startPage);
        if (startPage > 0) {
            ret.push('···');
        }
        for (var i = startPage; i < startPage + pageRange && i < length; i++) {
            var n = i + 1;
            ret.push(n);
        }
        // if (startPage < Math.floor(runtime.pagedItems.length / runtime.pageRange) * runtime.pageRange) {
        //     ret.push('···');
        // }
        //runtime.rangedPages = ret;
        return ret;
        //alert('Paging: ' + runtime.currentPage + "|" + startPage + "|" + runtime.rangedPages);
    };

    bMonitor.prevPage = function () {
        if (bMonitor.currentPage > 0) {
            bMonitor.setPage(bMonitor.currentPage);
        }
    };

    bMonitor.nextPage = function () {
        if (bMonitor.currentPage < bMonitor.pagedItems.length - 1) {
            bMonitor.setPage(bMonitor.currentPage + 2);
        }
    };

    bMonitor.setPage = function (page) {
        if (page == '···') {
            return;
        }
        //alert('page:' + page + '|' + this.n);
        bMonitor.currentPage = page - 1;
        bMonitor.range();
    };

    // functions have been describe process the data for display
    //bMonitor.search();

    // change sorting order
    bMonitor.sort_by = function (newSortingOrder) {
        if (sortingOrder == newSortingOrder)
            reverse = !reverse;

        sortingOrder = newSortingOrder;

        // icon setup
        $('th i').each(function () {
            // icon reset
            $(this).removeClass().addClass('icon-sort');
        });
        if (reverse)
            $('th.' + new_sorting_order + ' i').removeClass().addClass('icon-chevron-up');
        else
            $('th.' + new_sorting_order + ' i').removeClass().addClass('icon-chevron-down');
    };

    // ======= MODAL =======
    bMonitor.animationsEnabled = true;
    bMonitor.open = function (boiler, size, parentSelector) {
        currentBoiler = boiler;
        if (boiler.Fuel.Type.Id === 1 || boiler.Fuel.Type.Id === 4) {
            modalTemplate = '/directives/modal/boiler_calculate_coal.html';
        } else {
            modalTemplate = '/directives/modal/boiler_calculate_gas.html';
        }
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: bMonitor.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: modalTemplate,
            controller: 'ModalCalcCtl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            // resolve: {
            // items: function () {
            //     return confAlarm.items;
            // }
            // }
        });

        modalInstance.result.then(function (selectedItem) {
            // confAlarm.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    bMonitor.initBap = function () {
        // console.warn("bMonitor.initBap", bMonitor.filteredItems);
        if (!bMonitor.filteredItems || bMonitor.filteredItems.length <= 0) {
            console.warn("Boiler List IS Empty!");
            return;
        }
        if ( $('#map-container').size() === 0) {
            // console.warn("There IS NO #map-container!");
            return;
        }
        bMonitor.bMap = new BMap.Map("map-container"); // 创建地图实例
        bMonitor.bMap.addControl(new BMap.NavigationControl());
        bMonitor.bMap.addControl(new BMap.ScaleControl());
        bMonitor.bMap.addControl(new BMap.OverviewMapControl());
        bMonitor.bMap.enableContinuousZoom();

        var long = 0;
        var lat = 0;
        var count = 0;
        for (var i = 0; i < bMonitor.filteredItems.length; i++) {
            var b = bMonitor.filteredItems[i];
            if (!b.Address || b.Address.Longitude == 0 || b.Address.Latitude == 0) {
                continue;
            }
            var longitude = b.Address.Longitude;
            var latitude = b.Address.Latitude;
            count ++;
            long += longitude;
            lat += latitude;
            //console.warn(longitude, latitude, long, lat);
            var point = new BMap.Point(longitude, latitude);
            var marker = new BMap.Marker(point);
            var label = new BMap.Label(i);
            var offsetX = 0;
            var fontSize = 12;
            if (i < 10) {
                offsetX = 5;
            } else if (i < 100) {
                offsetX = 1;
            } else if (i < 1000) {
                fontSize = 10;
            }
            label.setStyle({
                'font-family': 'sans-serif',
                'font-size': fontSize + 'px',
                'text-align': 'center',
                'color': '#fff',
                'border': 'none',
                'background-color': 'transparent'});
            label.setOffset(new BMap.Size(offsetX, 2));
            marker.setTitle(b.Name);
            marker.setLabel(label);
            marker.addEventListener("click", function(){
                //alert("Clicked:" + b.Name);
                $location.hash('b' + b.num);
            });
            marker.addEventListener("dblclick", function(){
                $state.go("runtime.dashboard", {boiler: b.Uid});
            });
            bMonitor.bMap.addOverlay(marker);
            marker.setAnimation('BMAP_ANIMATION_DROP');
        }
        var cenLong = long / count;
        var cenLat = lat / count;
        console.warn("BMap Center", cenLong, cenLat, long, lat, count + "/" + bMonitor.filteredItems.length);
        var center = new BMap.Point(cenLong, cenLat);
        bMonitor.bMap.centerAndZoom(center, 10);  // 初始化地图，设置中心点坐标和地图级别
        // $scope.map.enableScrollWheelZoom(true);
        // 创建地址解析器实例
        $scope.myGeo = new BMap.Geocoder();
        /**
         * 监听地图点击事件，获取点击处建筑物名称
         */
        // $scope.map.addEventListener("click", function (e) {
        //     var pt = e.point;
        //     $scope.longitude = pt.lng;
        //     $scope.latitude = pt.lat;
        //     $scope.myGeo.getLocation(pt, function (rs) {
        //         var addComp = rs.addressComponents;
        //         /**
        //          * 将获取到的建筑名赋值给$scope.address
        //          */
        //         $scope.address = addComp.province != addComp.city ? addComp.province + addComp.city : addComp.city + addComp.district + addComp.street + addComp.streetNumber;
        //         /**
        //          * 通知angularjs更新视图
        //          */
        //         $scope.$digest();
        //     });
        // });
        /**
         * 初始化查询配置
         * @type {BMap.LocalSearch}
         */
        $scope.local = new BMap.LocalSearch(bMonitor.bMap, {
            renderOptions: {
                map: $scope.map,
                panel: "results",
                autoViewport: true,
                selectFirstResult: true
            },
            pageCapacity: 8
        });

        /**
         * 监听地址改变事件，当地址输入框的值改变时
         */
        $scope.$watch('address', function () {
            /**
             * 查询输入的地址并显示在地图上、调整地图视野
             */
            $scope.local.search($scope.address);
            /**
             * 将输入的地址解析为经纬度
             */
            $scope.myGeo.getPoint($scope.address, function (point) {
                if (point) {
                    /**
                     * 将地址解析为经纬度并赋值给$scope.longitude和$scope.latitude
                     */
                    $scope.longitude = point.lng;
                    $scope.latitude = point.lat;
                }
            });
        });
    };

    bMonitor.initAmChartDailyFlowAvg = function (dataProvider, graphs) {
        if ("undefined" != typeof AmCharts && 0 !== $("#dashboard_amchart_3d").size()) {

            var chart = AmCharts.makeChart("dashboard_amchart_3d", {
                theme: "light",
                type: "serial",
                language: "zh",
                dataDateFormat: 'MMM DDD',
                dataProvider: dataProvider,
                valueAxes: [{
                    stackType: "3d",
                    unit: "t",
                    position: "left",
                }],
                startDuration: 1,
                graphs: graphs,
                plotAreaFillAlphas: 0.1,
                depth3D: 120,
                angle: 60,
                categoryField: "num",
                categoryAxis: {
                    // minPeriod: "mm",
                    // parseDates: true,
                    // equalSpacing: true,
                    axisAlpha: 0.2,
                    gridPosition: "start"
                    //title: "本月平均每日流量"
                },
                export: {
                    enabled: true,
                }
            });
            // jQuery('.chart-input').off().on('input change',function() {
            //     var property	= jQuery(this).data('property');
            //     var target		= chart;
            //     chart.startDuration = 0;
            //
            //     if ( property == 'topRadius') {
            //         target = chart.graphs[0];
            //         if ( this.value == 0 ) {
            //             this.value = undefined;
            //         }
            //     }
            //
            //     target[property] = this.value;
            //     chart.validateNow();
            // });
        }
    };

    bMonitor.initAmChartHeat = function (dataProvider, graphs) {
        if ("undefined" != typeof AmCharts && 0 !== $("#dashboard_amchart_3d").size()) {

            var chart = AmCharts.makeChart("dashboard_amchart_3d", {
                theme: "light",
                type: "serial",
                language: "zh",
                dataDateFormat: 'MMM DDD',
                dataProvider: dataProvider,
                valueAxes: [{
                    stackType: "3d",
                    unit: "t",
                    position: "left",
                }],
                startDuration: 1,
                graphs: graphs,
                plotAreaFillAlphas: 0.1,
                depth3D: 120,
                angle: 60,
                categoryField: "num",
                categoryAxis: {
                    // minPeriod: "mm",
                    // parseDates: true,
                    // equalSpacing: true,
                    axisAlpha: 0.2,
                    gridPosition: "start"
                    //title: "本月平均每日流量"
                },
                export: {
                    enabled: true,
                }
            });
            // jQuery('.chart-input').off().on('input change',function() {
            //     var property	= jQuery(this).data('property');
            //     var target		= chart;
            //     chart.startDuration = 0;
            //
            //     if ( property == 'topRadius') {
            //         target = chart.graphs[0];
            //         if ( this.value == 0 ) {
            //             this.value = undefined;
            //         }
            //     }
            //
            //     target[property] = this.value;
            //     chart.validateNow();
            // });
        }
    };

    bMonitor.initAmChartPie = function () {
        if (!$rootScope.boilers || $rootScope.boilers.length === 0) {
            return;
        }

        if ("undefined" !== typeof AmCharts && 0 !== $("#dashboard_amchart_pie").size()) {
            var chartData = [];

            chartData.push({
                range: 'D＜1',
                count: 0
            });
            chartData.push({
                range: '1＜D≤2',
                count: 0
            });
            chartData.push({
                range: '2＜D≤8',
                count: 0
            });
            chartData.push({
                range: '8＜D≤20',
                count: 0
            });
            chartData.push({
                range: 'D＞20',
                count: 0
            });
            for (var i = 0; i < $rootScope.boilers.length; i++) {
                var boiler = $rootScope.boilers[i];
                var rate = boiler.EvaporatingCapacity;
                if (rate <= 1) {
                    chartData[0].count++;
                } else if (rate <= 2) {
                    chartData[1].count++;
                } else if (rate <= 8) {
                    chartData[2].count++;
                } else if (rate <= 20) {
                    chartData[3].count++;
                } else {
                    chartData[4].count++;
                }
            }

            var balloonText = function (graphDataItem, graph) {
                // console.error("balloonText", graphDataItem, graph);
                var title = graphDataItem.title;
                var count = graphDataItem.value;
                var percents = graphDataItem.percents.toFixed(2);
                return title + "<br><b>" + count + "台</b><br>(" + percents + "%)";
            };

            var chart = new AmCharts.AmPieChart();
            chart.theme = "light";
            chart.language = "zh";
            chart.valueField = "count";
            chart.titleField = "range";
            chart.colors= [
                "#84b761",
                "#fdd400",
                "#5fbfdb",
                "#c4e479",
                "#cd82ad",
                "#2f4074",
                "#448e4d",
                "#b7b83f",
                "#b9783f",
                "#b93e3d",
                "#913167"
            ];
            chart.startDuration = 1;

            chart.plotAreaFillAlphas = 0.1;
            chart.outlineAlpha = .4;
            chart.depth3D = 12;
            chart.angle = 30;
            chart.labelRadius = 16;
            chart.radius = 120;
            chart.legend={
                position:"bottom",
                marginRight:10,
                markerSize: 10,
                valueText: "",
                align: "center",
                autoMargins:false
            };
            chart.accessibleLabel = "[[title]]<br>[[value]] 台 ([[percents]]%)";
            chart.labelText = "[[title]]<br>[[percents]]%";
            chart.balloonFunction = balloonText;

            chart.export = { enabled: true };

            chart.dataProvider = chartData;

            // console.error("dataProvider:", chartData);

            chart.write("dashboard_amchart_pie");
        }
    };

});

angular.module('BoilerAdmin').controller('ModalCalcCtl', function ($uibModalInstance, $rootScope, $http) {
    var $modal = this;
    $modal.editing = false;
    $modal.boiler = currentBoiler;

    if (!currentBoiler) {
        $uibModalInstance.close('null object');
        return;
    }

    // console.warn("ModalCalcCtl:", $modal.boiler, $modal.boiler.Calculate);

    $modal.initCalc = function () {
        $modal.data = {};
        $modal.data.boiler_id = $modal.boiler.Uid;
        $modal.data.fuel_type_id = $modal.boiler.Fuel.Type.Id;

        $http.post('/boiler_runtime_instants/', {
            uid: $modal.boiler.Uid,
            runtimeQueue: [1014, 1021, 1016]
        }).then(function (res) {
            // console.warn("Get Instant Tempers:", res);
            $modal.data.smoke_temper = res.data[0].Value;
            $modal.data.wind_temper = res.data[1].Value;
            $modal.data.smoke_o2 = res.data[2].Value;
        }, function (err) {
            console.warn("Get Instant Tempers Error:", err);
        });

        $http.get('/boiler_calculate_parameter/?boiler=' + currentBoiler.Uid)
            .then(function (res) {
                currentBoiler.Calculate = res.data;

                initCalcParam(currentBoiler);
            });

        var initCalcParam = function (boiler) {
            if (!boiler.Calculate) {
                return;
            }

            var calcParam = boiler.Calculate;

            $modal.data.parameter_id = calcParam.Uid;
            if ($modal.data.fuel_type_id === 1 || $modal.data.fuel_type_id === 4) {
                // COAL
                $modal.data.qnetvar = calcParam.CoalQnetvar;
                $modal.data.aar = calcParam.CoalAar;
                $modal.data.mar = calcParam.CoalMar;
                $modal.data.vdaf = calcParam.CoalVdaf;
                $modal.data.clz = calcParam.CoalClz;
                $modal.data.clm = calcParam.CoalClm;
                $modal.data.cfh = calcParam.CoalCfh;
                $modal.data.ded = calcParam.CoalDed;
                $modal.data.dsc = calcParam.CoalDsc;
                $modal.data.alz = calcParam.CoalAlz;
                $modal.data.alm = calcParam.CoalAlm;
                $modal.data.afh = calcParam.CoalAfh;
                $modal.data.tlz = calcParam.CoalTlz;
                $modal.data.ct_lz = calcParam.CoalCtLz;

                // $modal.data.apy = 0;
                // $modal.data.q2 = 0;
                // $modal.data.q4 = 0;
                // $modal.data.q6 = 0;

                $modal.data.m = calcParam.CoalM;
                $modal.data.n = calcParam.CoalN;

                $modal.data.q3 = calcParam.CoalQ3;
                $modal.data.q5 = calcParam.ConfParam1;
            } else {
                // GAS
                $modal.data.ded = calcParam.GasDed;

                $modal.data.m = calcParam.GasM;
                $modal.data.n = calcParam.GasN;

                $modal.data.q3 = calcParam.GasQ3;
                $modal.data.q5 = calcParam.ConfParam1;
            }
        }
    };

    $modal.reset = function () {
        $modal.initCalc();
    };

    $modal.calculate = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();

        $http.post("/boiler_calculate/", $modal.data)
            .then(function (res) {
                console.warn("Boiler Calculate Res:", res);
                $modal.data.q2 = res.data.q2;
                $modal.data.q3 = res.data.q3;
                $modal.data.q4 = res.data.q4;
                $modal.data.q5 = res.data.q5;
                $modal.data.q6 = res.data.q6;
                $modal.data.excessAir = res.data.apy;
                $modal.data.heat = res.data.Heat;
            }, function (err) {
                swal({
                    title: "参数计算失败",
                    text: err.data,
                    type: "error"
                });
            });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.save = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();
        $http.post("/boiler_calculate_parameter_update/", $modal.data)
            .then(function (res) {
            console.log("res", res);
            swal({
                title: "计算参数更新成功",
                type: "success"
            }).then(function () {
                $rootScope.getBoilerCalculateParameter($modal.boiler);
            });
        }, function (err) {
            swal({
                title: "计算参数更新失败",
                text: err.data,
                type: "error"
            });
        });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentBoiler = null;
    };

    $modal.initCalc();
});

// Please note that the close and dismiss bindings are from $uibModalInstance.

angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: modalTemplate,
    bindings: {
        resolve: '<',
        close: '&',
        dismiss: '&'
    },
    controller: function () {
        var $ctrl = this;

        $ctrl.$onInit = function () {
            // $ctrl.items = $ctrl.resolve.items;
            $ctrl.selected = {
                // item: $ctrl.items[0]
            };
        };

        $ctrl.ok = function () {
            $ctrl.close({$value: $ctrl.selected.item});
        };

        $ctrl.cancel = function () {
            $ctrl.dismiss({$value: 'cancel'});
        };
    }
});

var bMonitor;

var currentBoiler;
var modalTemplate = '/directives/modal/boiler_calculate_gas.html';