angular.module('BoilerAdmin')
    .factory('Excel',function($window){
        var uri = 'data:application/vnd.ms-excel;charset=utf-8;base64,',
            template = '<html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40"><head><!--[if gte mso 9]><xml><x:ExcelWorkbook><x:ExcelWorksheets><x:ExcelWorksheet><x:Name>{worksheet}</x:Name><x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions></x:ExcelWorksheet></x:ExcelWorksheets></x:ExcelWorkbook></xml><![endif]--></head><body><table>{table}</table></body></html>',
            base64 = function(s){return $window.btoa(unescape(encodeURIComponent(s)));},
            format = function(s,c){return s.replace(/{(\w+)}/g,function(m,p){return c[p];})};
        return {
            tableToExcel:function(tableId, worksheetName){
                var table = $(tableId),
                    ctx = {
                        worksheet: worksheetName,
                        table: table.html()
                    },
                    href = uri + base64(format(template, ctx));
                return href;
            }
        };
    })
    .controller('BoilerHistoryController', function($rootScope, $scope, $http, $location, $timeout, $log, $document, Excel, moment, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
        bHistory = this;

        bHistory.isDone = false;
        bHistory.isEmpty = false;

        $scope.$on('$viewContentLoaded', function() {
            // initialize core components

            // $log.info("init Boiler History Controller!");
            // createGauges();
            // setInterval(updateGauges, 5000);
            App.initAjax();

            // set sidebar closed and body solid layout mode
            $rootScope.settings.layout.pageContentWhite = true;
            $rootScope.settings.layout.pageBodySolid = true;
            $rootScope.settings.layout.pageSidebarClosed = false;
        });

        $scope.exportToExcel = function(tableId) { //ex: '#my-table'
            $scope.exportHref = Excel.tableToExcel(tableId, 'sheet name');
            $timeout(function(){
                location.href = $scope.exportHref;
                }, 100); // trigger download
        };

        bHistory.excelExport = function () {
            var excelData = [];
            var excelName = $rootScope.boiler.Name;
            var start = moment(bHistory.startDate).format('YYYY.MM.DD');
            var end = moment(bHistory.endDate).format('YYYY.MM.DD');
            excelName += " (" + start + " ~ " + end + ")";
            var xdp = {
                0: "采样时间"
            };
            for (var p = 0; p < bHistory.parameters.length; p++) {
                var param = bHistory.parameters[p];
                xdp[param.Id] = param.Name + " " + param.Unit;
            }
            excelData.push(xdp);

            for (var i = 0; i < bHistory.datasource.length; i++) {
                var d = bHistory.datasource[i];
                var xd = {
                    0: moment(d.date).format('YYYY-MM-DD HH:mm')
                };

                for (var j = 0; j < Object.keys(d).length; j++) {
                    var k = Object.keys(d)[j];
                    if (parseInt(k) > 1000) {
                        xd[k] = d[k].value;
                    }
                }

                excelData.push(xd);
            }

            var res = alasql('SELECT * INTO XLSX("' + excelName + '.xlsx", {headers:true}) FROM ?', [excelData]);
            console.log(res);
        };

        bHistory.dtOptions = DTOptionsBuilder.newOptions()
            .withPaginationType('full_numbers');
            //.withOption('rowCallback', rowCallbackHistory);

        bHistory.datasource = [];

        bHistory.refreshDataTables = function () {
            // $log.info("history.refreshDataTables!");
            var p = $location.search();
            bHistory.pids = [];
            $http.post('/boiler_runtime_history/', {
                uid: p['boiler'],
                runtimeQueue: bHistory.pids,
                startDate: bHistory.startDate,
                endDate: bHistory.endDate,
                limit: -1
            }).then(function (res) {
                $log.warn("Runtime History Resp1:", res);

                bHistory.datasource = [];
                bHistory.parameters = res.data.parameter;

                bHistory.pids = [];
                for (var i = 0; i < bHistory.parameters.length; i++) {
                    var param = bHistory.parameters[i];
                    bHistory.pids.push(param.Id);
                }

                Ladda.create(document.getElementById('history_range_today')).stop();
                Ladda.create(document.getElementById('history_range_week')).stop();
                Ladda.create(document.getElementById('history_range_month')).stop();
                if (!res.data.history || res.data.history.length === 0) {
                    bHistory.isDone = true;
                    bHistory.isEmpty = true;
                    return;
                }

                for (var i = 0; i < res.data.history.length; i++) {
                    var hData = res.data.history[i];
                    var d = {};
                    d.num = i;
                    d.id = i;
                    d.date = new Date(hData.date);

                    for (var k = 0; k < res.data.parameter.length; k++) {
                        var ap = res.data.parameter[k];
                        switch (ap.Id) {
                            case 1021:
                                ap.Name = "环境温度";
                                break;
                            case 1202:
                                ap.Name = "过量空气系数";
                                break;
                        }

                        var key = ap.Id.toString();
                        d[key] = {
                            value: '-',
                            alarm: -1
                        };
                        for (var j = 0; j < hData.data.length; j++) {
                            var h = hData.data[j];
                            if (h.pid === ap.Id) {
                                d[key].value = h.val;
                                d[key].alarm = h.alm;
                            }
                        }

                    }

                    bHistory.datasource.push(d);
                }

                $scope.groupToPages();
                $scope.setPage(1);

                bHistory.isDone = true;
                bHistory.isEmpty = false;
                $log.warn("History Data1:", bHistory);
            });
        };

        bHistory.setDataRange = function (range) {
            var startDate = new Date();
            var endDate = new Date();
            Ladda.create(document.getElementById('history_range_today')).stop();
            Ladda.create(document.getElementById('history_range_week')).stop();
            Ladda.create(document.getElementById('history_range_month')).stop();
            Ladda.create(document.getElementById('history_range_' + range)).start();
            switch (range) {
                case 'today':
                    startDate.setHours(0);
                    startDate.setMinutes(0);
                    break;

                case 'week':
                    startDate.setDate(startDate.getDate() - 7);
                    startDate.setHours(0);
                    startDate.setMinutes(0);
                    break;

                case 'month':
                    startDate.setDate(1);
                    startDate.setHours(0);
                    startDate.setMinutes(0);
                    break;

                default:
                    break;
            }

            // endDate.setMinutes(0);
            bHistory.startDate = startDate;
            bHistory.endDate = endDate;

            bHistory.dataRange = range;

            bHistory.refreshDataTables();
        };

        bHistory.dateChanged = function () {
            console.warn("bHistory.dateChanged:", bHistory.startDate, "-", bHistory.endDate);
            if (bHistory.startDate < bHistory.endDate) {
                bHistory.refreshDataTables();
            } else {
                bHistory.datasource = [];
            }
        };

        var sortingOrder = '';
        var reverse = false;

        var groupedItems = [];
        var itemsPerPage = 50;
        var pageRange = 100;

        $scope.pagedItems = [];
        $scope.currentPage = 0;
        $scope.filterLen = 0;
        //bMonitor.rangedPages = [];

        $scope.matchNum = 0;

        // calculate page in place
        $scope.groupToPages = function () {
            $scope.pagedItems = [];
            for (var i = 0; i < bHistory.datasource.length; i++) {
                if (i % itemsPerPage === 0) {
                    $scope.pagedItems[Math.floor(i / itemsPerPage)] = [bHistory.datasource[i]];
                } else {
                    $scope.pagedItems[Math.floor(i / itemsPerPage)].push(bHistory.datasource[i]);
                }
            }

            // console.warn("$scope.groupToPages():", $scope.pagedItems);
        };

        $scope.range = function () {
            var ret = [];
            var length = $scope.pagedItems.length;
            var startPage = Math.floor($scope.currentPage / pageRange) * pageRange;
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

        $scope.prevPage = function () {
            if ($scope.currentPage > 0) {
                $scope.setPage($scope.currentPage);
            }
        };

        $scope.nextPage = function () {
            if ($scope.currentPage < $scope.pagedItems.length - 1) {
                $scope.setPage($scope.currentPage + 2);
            }
        };

        $scope.setPage = function (page) {
            if (page === '···') {
                return;
            }
            //alert('page:' + page + '|' + this.n);
            $scope.currentPage = page - 1;
            $scope.range();
        };

        // functions have been describe process the data for display
        //bMonitor.search();

        // change sorting order
        $scope.sort_by = function (newSortingOrder) {
            if (sortingOrder === newSortingOrder)
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

        $scope.mytime = new Date();
        $scope.mytime.setMinutes(0);

        $scope.hstep = 1;
        $scope.mstep = 30;

        $scope.ismeridian = true;
        $scope.toggleMode = function() {
            $scope.ismeridian = ! $scope.ismeridian;
        };

        $scope.update = function() {
            var d = new Date();
            d.setHours( 14 );
            d.setMinutes( 0 );
            $scope.mytime = d;
        };

        $scope.changed = function () {
            $log.log('Time changed to: ' + $scope.mytime);
        };

        $scope.clear = function() {
            $scope.mytime = null;
        };

        $scope.today = function() {
            $scope.dt = new Date();
        };
        $scope.today();

        $scope.clear = function() {
            $scope.dt = null;
        };

        $scope.inlineOptions = {
            customClass: getDayClass,
            minDate: new Date(),
            showWeeks: true
        };

        $scope.dateMinOptions = {
            // dateDisabled: disabled,
            formatYear: 'yy',
            maxDate: bHistory.endDate,
            minDate: new Date(2017, 5, 1),
            startingDay: 0
        };

        $scope.dateMaxOptions = {
            // dateDisabled: disabled,
            formatYear: 'yy',
            maxDate: new Date(),
            minDate: bHistory.startDate,
            startingDay: 0
        };

        // Disable weekend selection
        function disabled(data) {
            var date = data.date,
                mode = data.mode;
            return mode === 'day' && (date.getDay() === 0 || date.getDay() === 6);
        }

        // $scope.toggleMin = function() {
        //     $scope.inlineOptions.minDate = $scope.inlineOptions.minDate ? null : new Date();
        //     $scope.dateOptions.minDate = $scope.inlineOptions.minDate;
        // };

        // $scope.toggleMax = function() {
        //     $scope.inlineOptions.minDate = $scope.inlineOptions.minDate ? null : new Date();
        //     $scope.dateOptions.minDate = $scope.inlineOptions.minDate;
        // };

        // $scope.toggleMin();

        $scope.open1 = function() {
            $scope.popup1.opened = true;
        };

        $scope.open2 = function() {
            $scope.popup2.opened = true;
        };

        $scope.setDate = function(year, month, day) {
            $scope.dt = new Date(year, month, day);
        };

        $scope.popup1 = {
            opened: false
        };

        $scope.popup2 = {
            opened: false
        };

        var tomorrow = new Date();
        tomorrow.setDate(tomorrow.getDate() + 1);
        var afterTomorrow = new Date();
        afterTomorrow.setDate(tomorrow.getDate() + 1);
        $scope.events = [
            {
                date: tomorrow,
                status: 'full'
            },
            {
                date: afterTomorrow,
                status: 'partially'
            }
        ];

        function getDayClass(data) {
            var date = data.date,
                mode = data.mode;
            if (mode === 'day') {
                var dayToCheck = new Date(date).setHours(0,0,0,0);

                for (var i = 0; i < $scope.events.length; i++) {
                    var currentDay = new Date($scope.events[i].date).setHours(0,0,0,0);

                    if (dayToCheck === currentDay) {
                        return $scope.events[i].status;
                    }
                }
            }

            return '';
        }
    });

var bHistory;