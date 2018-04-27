angular.module('BoilerAdmin').controller('BoilerInfoController', function($rootScope, $scope, $http, $timeout, $location, $log, $window, $state, $uibModal, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    bInfo = this;
    bInfo.isDone = false;

    bInfo.dataMode = "edit";

    $rootScope.$on('$viewContentLoading', function(event, viewConfig){
        // Access to all the view config properties.
        // and one special property 'targetView'
        // viewConfig.targetView
    });

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        // bInfo.refreshDataTables();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    bInfo.new = function () {
        bInfo.currentData = null;
        bInfo.openBasic('lg');
    };

    bInfo.initEditing = function () {
        bInfo.editing = {
            basic: false,
            address: false,
            maintain: false,
            terminal: false,
            calc: false
        }
    };

    bInfo.reset = function () {
        bInfo.basic = {
            Name: bInfo.currentData.Name,
            RegisterCode: bInfo.currentData.RegisterCode,
            DeviceCode: bInfo.currentData.DeviceCode,
            FactoryNumber: bInfo.currentData.FactoryNumber,
            ModelCode: bInfo.currentData.ModelCode,
            CertificateNumber: bInfo.currentData.CertificateNumber,

            Usage: bInfo.currentData.Usage,
            Medium: bInfo.currentData.Medium,
            Fuel: bInfo.currentData.Fuel,
            Form: bInfo.currentData.Form,
            EvaporatingCapacity: bInfo.currentData.EvaporatingCapacity,

            RegisterOrg: bInfo.currentData.RegisterOrg,
            Enterprise: bInfo.currentData.Enterprise,
            Factory: bInfo.currentData.Factory,
            Maintainer: bInfo.currentData.Maintainer,
            Supervisor: bInfo.currentData.Supervisor,

            Links: []
        };

        for (var i in bInfo.currentData.OrganizationsLinked) {
            var uid = bInfo.currentData.OrganizationsLinked[i].Uid;
            var link = {};
            link.num = ++i;
            link.uid = uid;

            link.name = "";
            link.type = 0;
            link.typeName = "";

            bInfo.basic.Links.push(link);
        }

        bInfo.currentData.Links = bInfo.basic.Links;

        bInfo.updateCurrentLinks();
        bInfo.updateCurrentLinksType();

        bInfo.initBap();

        bInfo.initEditing();
    };

    bInfo.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
        //.withOption('rowCallback', rowBoilerCallback);

    bInfo.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6).notSortable()
    ];

    var someClickHandler = function(info) {
        bInfo.message = info.Uid + ' - ' + info.Name;

        if (!info.Calculate) {
            $http.get('/boiler_calculate_parameter/?boiler=' + info.Uid)
                .then(function (res) {
                    info.Calculate = res.data;
                    if (info.Calculate) {
                        info.calc = info.Calculate;
                    }
                });
        }

        $state.go("boiler.info", {boiler: info.Uid, from: 'boiler-list'});
    };

    /*
    function rowBoilerCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        $('td', nRow).unbind('click');
        $('td', nRow).bind('click', function() {
            bInfo.editing = false;
            bInfo.currentData = bInfo.datasource[aData[0] - 1];
            bInfo.currentData.aName = bInfo.currentData.Name;

            $scope.$apply(function() {
                someClickHandler(bInfo.currentData);
            });
        });
        return nRow;
    }
    */

    bInfo.viewInfo = function (data) {
        console.warn("Boiler View Uid:", data);

        bInfo.currentData = data;

        bInfo.setMode('edit');

        //$state.go("boiler.info", {boiler: info.Uid, from: 'boiler-list'});
        someClickHandler(bInfo.currentData);
    };

    bInfo.setMode = function (mode) {
        bInfo.dataMode = mode;
    };

    bInfo.refreshDataTables = function () {
        var p = $location.search()['boiler'];

        var datasource = $rootScope.boilers;

        var num = 0;
        angular.forEach(datasource, function (d, key) {
            d.num = ++num;

            $http.get('/boiler/state/is_burning/?boiler=' + d.Uid)
                .then(function (res) {
                    // console.error("Fetch Status Resp:", res.data, boiler.Name);
                    d.isBurning = (res.data.value === "true");
                }, function (err) {
                    console.error('Fetch Status Err!', err);
                });

            if (d.TerminalsCombined) {
                for (var i = 0; i < d.TerminalsCombined.length; i++) {
                    var terminal = d.TerminalsCombined[i];

                    terminal.tid = terminal.TerminalCode.toString();
                    if (terminal.tid.length < 6) {
                        for (var l = terminal.tid.length; l < 6; l++) {
                            terminal.tid = "0" + terminal.tid;
                        }
                    }

                    terminal.online = terminal.IsOnline ? "在线" : "离线";
                    terminal.setId = terminal.Remark;

                }

            }

            if (d.Calculate) {
                d.calc = d.Calculate;
            }

            if (p && p.length > 0 && d.Uid === p) {
                bInfo.currentData = d;

                bInfo.reset();
            }
            // console.log(d);
        });

        bInfo.datasource = datasource;
        bInfo.isDone = true;
        setTimeout(function () {
            App.stopPageLoading();
        }, 800);
    };

    bInfo.updateCurrentLinks = function () {
        for (var i in bInfo.currentData.Links) {
            var link = bInfo.currentData.Links[i];
            for (var j in $rootScope.organizations) {
                var og = $rootScope.organizations[j];
                if (og.Uid === link.uid) {
                    link.name = og.name;
                    link.type = og.typeId;
                    break;
                }
            }
        }
    };

    bInfo.updateCurrentLinksType = function () {
        for (var i in bInfo.currentData.Links) {
            var link = bInfo.currentData.Links[i];
            if (!link.type || link.type <= 0) {
                continue;
            }
            for (var j in $rootScope.organizationTypes) {
                var tp = $rootScope.organizationTypes[j];
                // console.warn("Linke:", tp, link.type);
                if (tp.id === link.type) {
                    // console.error("Get Linked:", link.uid, tp.name);
                    link.typeName = tp.name;
                    break;
                }
            }

            bInfo.currentData.Links[i] = link;
        }
    };

    $rootScope.$watch('boilers', function () {
        // console.error("$rootScope.$watch('boilers')");
        bInfo.refreshDataTables();
    });

    $rootScope.$watch('organizations', function () {
        if (!$rootScope.organizations) {
            return;
        }

        bInfo.updateCurrentLinks();
        bInfo.updateCurrentLinksType();
    });

    $rootScope.$watch('organizationTypes', function () {
        if (!$rootScope.organizationTypes) {
            return;
        }

        bInfo.updateCurrentLinksType();
    });

    bInfo.initBap = function () {
        if (!bInfo.currentData) {
            console.warn("Current Boiler List IS Empty!");
            return;
        }
        if ( $('#map-container').size() === 0) {
            console.warn("There IS NO #map-container!");
            return;
        }
        bInfo.bMap = new BMap.Map("map-container"); // 创建地图实例
        //bInfo.bMap.addControl(new BMap.NavigationControl());
        bInfo.bMap.addControl(new BMap.ScaleControl());
        bInfo.bMap.addControl(new BMap.OverviewMapControl());

        var b = bInfo.currentData;

        var longitude = 0;
        var latitude = 0;

        var geolocation = new BMap.Geolocation();
        geolocation.getCurrentPosition(
            //获取位置信息成功
            function(position){
                if(this.getStatus() == BMAP_STATUS_SUCCESS){
                    console.warn('您的位置：' + position.point.lng + ',' + position.point.lat);
                    longitude = position.point.lng;
                    latitude  = position.point.lat;
                    // 根据坐标得到地址描述
                } else {
                    console.error('无法获取定位信息');
                }
            },{
                // 指示浏览器获取高精度的位置，默认为false
                enableHighAccuracy: true,
                // 指定获取地理位置的超时时间，默认不限时，单位为毫秒
                // timeout: 5000,
                // 最长有效期(30S)，在重复获取地理位置时，此参数指定多久再次获取位置
                maximumAge: 30*1000
            });

        var point = new BMap.Point(longitude, latitude);

        if (b.Address && b.Address.Longitude !== 0 && b.Address.Latitude !== 0) {
            console.warn("There IS Address!");
            longitude = b.Address.Longitude;
            latitude = b.Address.Latitude;

            point = new BMap.Point(longitude, latitude);
            var marker = new BMap.Marker(point);

            marker.setTitle(b.Name);
            marker.setAnimation('BMAP_ANIMATION_DROP');

            bInfo.bMap.addOverlay(marker);
            // 初始化地图，设置中心点坐标和地图级别
        }

        bInfo.bMap.centerAndZoom(point, 12);
    };

    bInfo.bindSet = function () {
        if (!bInfo.currentData) {
            console.warn("There is no device to bind!");
            return;
        }

        bInfo.openBind();
    };

    bInfo.unbind = function (terminal) {
        swal({
            title: "解除终端#" + bInfo.currentData.Name + "\n与该设备的绑定?",
            text: "解除绑定后，将无法收到来自 " + terminal.Name + " 的运行时数据。",
            type: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            cancelButtonColor: '#3085d6',
            confirmButtonText: '解绑',
            cancelButtonText: '取消'
        }).then(function () {
            $http.post("/boiler_unbind/", {
                boiler_id: bInfo.currentData.Uid,
                terminal_id: terminal.Uid
            }).then(function (res) {
                swal({
                    title: "绑定已解除",
                    text: "该锅炉已不再接收 " + terminal.Name + " 相关数据，如需重新接入，请通过终端绑定流程进行再次绑定。",
                    type: "success"
                });
                $rootScope.getBoilerList();
            }, function (err) {
                swal({
                    title: "解除绑定失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

    bInfo.openBasic = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/boiler_info_basic.html',
            controller: 'ModalBoilerInfoBasicCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                currentData: function () {
                    //console.info("resolve:", bInfo.currentData);
                    return bInfo.currentData;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            // terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    bInfo.openLocation = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/boiler_info_location.html',
            controller: 'ModalBoilerInfoLocationCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                currentData: function () {
                    //console.info("resolve:", bInfo.currentData);
                    return bInfo.currentData;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            // terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    bInfo.openMaintain = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/boiler_info_maintain.html',
            controller: 'ModalBoilerInfoMaintainCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                currentData: function () {
                    //console.info("resolve:", bInfo.currentData);
                    return bInfo.currentData;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            // terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    bInfo.openBind = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/boiler_bind.html',
            controller: 'ModalBoilerBindCtrl',
            controllerAs: '$modalBind',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex_sub',
            resolve: {
                currentData: function () {
                    //console.info("resolve:", bInfo.currentData);
                    return bInfo.currentData;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            //terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    bInfo.back = function () {
        $log.info('Ready to bInfo.back', $location.search()['from']);
        switch ($location.search()['from']) {
            case 'runtime':
                $state.go('runtime.stats', {boiler: bRuntime.boiler.Uid});
                break;
            case 'boiler-list':
                $state.go('boiler.dashboard');
                break;
            case 'terminal':
                $state.go('terminal');
                break;
            default:
                break;
        }
    };

    bInfo.edit = function (field) {
        switch (field) {
            case "basic":
                bInfo.openBasic('lg');
                break;
            case "location":
                bInfo.openLocation('lg');
                break;
            case "maintain":
                bInfo.openMaintain('lg');
                break;
        }
        //bInfo.editing[field] = true;
    };

    bInfo.delete = function () {
        if (!bInfo.currentData) {
            return;
        }

        swal.queue([{
            title: '确认删除该设备？',
            text: "注意：删除后将无法恢复，且将丢失该设备的所有运行时数据",
            type: "warning",

            showCancelButton: true,
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            showLoaderOnConfirm: true,
            closeOnConfirm: false,
            preConfirm: function () {
                return new Promise(function (resolve) {
                    Ladda.create(document.getElementById('boiler_delete')).start();
                    $http.post("/boiler_delete/", {
                        uid: bInfo.currentData.Uid
                    }).then(function (res) {
                        // console.error("boiler delete res:", res);
                        swal({
                            title: "设备删除成功",
                            type: "success"
                        }).then(function () {
                            resolve();
                        });

                        for (var i = 0; i < $rootScope.boilers.length; i++) {
                            var b = $rootScope.boilers[i];
                            if (b.Uid === bInfo.currentData.Uid) {
                                $rootScope.boilers.splice(i, 1);
                                break;
                            }
                        }

                        bInfo.refreshDataTables();
                        bInfo.back();

                        Ladda.create(document.getElementById('boiler_delete')).stop();
                    }, function (err) {
                        // console.error("boiler delete err:", err);
                        swal({
                            title: "设备删除失败",
                            text: err.data,
                            type: "error"
                        });
                        Ladda.create(document.getElementById('boiler_delete')).stop();
                    });
                })
            }
        }]);
    };
});

var bInfo;

angular.module('BoilerAdmin').controller('ModalBoilerInfoBasicCtrl', function ($uibModalInstance, $rootScope, $http, $log, currentData) {
    var $modal = this;
    //console.warn("CurrentData In Modal:", currentData);
    $modal.currentData = currentData;
    $modal.editingCode = true;

    $modal.links = [];

    $modal.initData = function (currentData) {
        if (currentData) {
            $modal.editingCode = false;

            $modal.data = {
                uid: currentData.Uid,
                name: currentData.Name,
                registerCode: currentData.RegisterCode,
                deviceCode: currentData.DeviceCode,
                factoryNumber: currentData.FactoryNumber,
                modelCode: currentData.ModelCode,
                certificateNumber: currentData.CertificateNumber,

                Usage: "工业锅炉",
                mediumId: currentData.Medium ? currentData.Medium.Id : -1,
                fuelId: currentData.Fuel ? currentData.Fuel.Uid : "",
                formId: currentData.Form ? currentData.Form.Id : -1,

                evaporatingCapacity: currentData.EvaporatingCapacity,

                RegisterOrg: currentData.RegisterOrg ? currentData.RegisterOrg : null,
                enterpriseId: currentData.Enterprise ? currentData.Enterprise.Uid : "",
                factoryId: currentData.Factory ? currentData.Factory.Uid : "",
                maintainerId: currentData.Maintainer ? currentData.Maintainer.Uid : "",
                supervisorId: currentData.Supervisor ? currentData.Supervisor.Uid : ""
            };

            for (var i in currentData.OrganizationsLinked) {
                var uid = currentData.OrganizationsLinked[i].Uid;
                var link = {};
                link.num = ++i;
                link.uid = uid;

                for (var j in $rootScope.organizations) {
                    var og = $rootScope.organizations[j];
                    if (og.Uid === uid) {
                        console.error("Get Linked:", uid, og.name);
                        link.name = og.name;
                        link.type = og.typeId;
                        break;
                    }
                }

                var orgs = [];
                for (var k in $rootScope.organizations) {
                    var og = $rootScope.organizations[k];
                    if (og.typeId === link.type) {
                        orgs.push(og);
                    }
                }

                link.orgs = orgs;

                $modal.links.push(link);
            }

        } else {
            $modal.data = {
                uid: "",
                name: "",
                registerCode: "",
                deviceCode: "",
                factoryNumber: "",
                modelCode: "",
                certificateNumber: "",

                Usage: "工业锅炉",
                mediumId: -1,
                fuelId: "",
                formId:  -1,

                evaporatingCapacity: 0,

                RegisterOrg: null,
                enterpriseId: "",
                factoryId: "",
                installedId: "",

                links: []
            }
        }
    };

    $modal.init = function () {
        $modal.mediums = [{ Id: -1, Name: '请选择...' }];
        $modal.forms = [{ Id: -1, Name: '请选择...' }];
        $modal.fuels = [{ Uid: '', Name: '请选择...' }];

        $modal.enterprises = [{ Uid: '', name: '请选择...' }];
        $modal.factories = [{ Uid: '', name: '请选择...' }];
        $modal.maintainers = [{ Uid: '', name: '请选择...' }];
        $modal.supervisors = [{ Uid: '', name: '请选择...' }];

        $modal.orgTypes = [];

        for (var i in $rootScope.organizations) {
            var org = $rootScope.organizations[i];
            switch (org.Type.TypeId) {
                case 2:
                    if ($modal.enterprises.indexOf(org) === -1) {
                        $modal.enterprises.push(org);
                    }
                    break;
                case 1:
                    if ($modal.factories.indexOf(org) === -1) {
                        $modal.factories.push(org);
                    }
                    break;
                case 5:
                    if ($modal.maintainers.indexOf(org) === -1) {
                        $modal.maintainers.push(org);
                    }
                    break;
                case 7:
                    if ($modal.supervisors.indexOf(org) === -1) {
                        $modal.supervisors.push(org);
                    }
                    break;
            }
        }

        for (var i in $rootScope.organizationTypes) {
            var t = $rootScope.organizationTypes[i];
            switch (t.id) {
                case 3:
                case 4:
                case 6:
                    $modal.orgTypes.push(t);
                    break;
            }
        }

        for (var i in $rootScope.boilerMediums) {
            var med = $rootScope.boilerMediums[i];
            if (med.Id === 0 || $modal.mediums.indexOf(med) > -1) {
                continue;
            }

            $modal.mediums.push(med);
        }

        for (var i in $rootScope.boilerForms) {
            var form = $rootScope.boilerForms[i];
            if (form.Id === 0 || $modal.forms.indexOf(form) > -1) {
                continue;
            }

            $modal.forms.push(form);
        }

        for (var i in $rootScope.fuels) {
            var fuel = $rootScope.fuels[i];
            if (fuel.Type.Id === 0 || fuel.Type.Id >= 5 || $modal.fuels.indexOf(fuel) > -1) {
                continue;
            }
            $modal.fuels.push(fuel);
        }

        $modal.initData(currentData);
    };

    $modal.init();

    $modal.addNewLink = function () {
        if ($modal.links.length >= 4) {
            return;
        }

        $modal.links.push({
            num: $modal.links.length,
            // type: -1,
            // uid: "",
            // name: ""
        });
    };

    $modal.removeLink = function (link) {
        for (var i in $modal.links) {
            if (link === $modal.links[i]) {
                $modal.links.splice(i, 1);
            }
        }
    };

    $modal.linkTypeChanged = function (link) {
        var orgs = [];
        for (var i in $rootScope.organizations) {
            var og = $rootScope.organizations[i];
            if (og.typeId === link.type) {
                orgs.push(og);
            }
        }

        link.orgs = orgs;
        link.uid = undefined;
    };

    $modal.save = function () {
        console.info("ready to update bInfo!");
        Ladda.create(document.getElementById('boiler_basic_submit')).start();
        $modal.data.links = [];
        for (var i in $modal.links) {
            var li = $modal.links[i];
            $modal.data.links.push({
                num: li.num,
                type: li.type,
                uid: li.uid
            });
        }
        $http.post("/boiler_update/?scope=basic", $modal.data)
            .then(function (res) {
                console.error("Update bInfo Resp:", res);
                $rootScope.getBoilerList();
                swal({
                    title: "设备基本信息更新成功",
                    type: "success"
                }).then(function () {

                });

                $uibModalInstance.close('success');
                // $modal.initData(res.data);
                // bInfo.currentData = res.data;
                // bInfo.reset();
                currentData = null;
                Ladda.create(document.getElementById('boiler_basic_submit')).stop();
            }, function (err) {
                swal({
                    title: "设备基本信息更新失败",
                    text: err.data,
                    type: "error"
                });
                Ladda.create(document.getElementById('boiler_basic_submit')).stop();
            });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        //currentData = null;
    };
});

angular.module('BoilerAdmin').controller('ModalBoilerInfoLocationCtrl', function ($uibModalInstance,  $rootScope, $scope, $filter, $http, $log, currentData) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.data = {};

    var needSearch = false;
    var needLocation = false;

    $modal.initData = function () {
        $scope.aProvince = $rootScope.locations[0];
        $scope.aProvince.Name = "所在区域";

        if (currentData) {
            $scope.uid = currentData.Uid;
            $scope.name = currentData.Name;

            $scope.address = currentData.Address.Address;
            $scope.location = currentData.Address.Location;
            $scope.longitude = currentData.Address.Longitude;
            $scope.latitude = currentData.Address.Latitude;
            var locationId = $scope.location.LocationId;
            getLocation(locationId, $rootScope.locations, "province");
        } else {
            $scope.uid = "";
            $scope.name = "";
            $scope.address = "";
            $scope.location = $modal.aProvince;
            $scope.longitude = 0.0;
            $scope.latitude = 0.0;
        }
    };

    $modal.changeProvince = function () {
        if (!$scope.aProvince || !$scope.aProvince.LocationId === 0) {
            return;
        }
        needLocation = true;
        $scope.location = $scope.aProvince;
    };

    $modal.changeCity = function () {
        if (!$scope.aCity) {
            return;
        }
        needLocation = true;
        $scope.location = $scope.aCity;
    };

    $modal.changeRegion = function () {
        if (!$scope.aRegion) {
            return;
        }
        needLocation = true;
        $scope.location = $scope.aRegion;

    };

    var getLocation = function (locationId, locationList, locationScope) {
        // console.warn("Init getLocation:", locationId, locationScope, locationList);
        for (var pi = 0; pi < locationList.length; pi++) {
            var local = locationList[pi];
            // console.warn("Ready to Find Local:", locationScope, local, locationId);
            if (local.LocationId === Math.floor(locationId / 10000) ||
                local.LocationId === Math.floor(locationId / 100) ||
                local.LocationId === locationId) {
                // console.error("Get Local:", locationScope, local, locationId);
                switch (locationScope) {
                    case "province":
                        $scope.aProvince = local;
                        break;
                    case "city":
                        $scope.aCity = local;
                        break;
                    case "region":
                        $scope.aRegion = local;
                        break;
                }
                break;
            }
        }

        if (locationId < 100) {
            return;
        }

        switch (locationScope) {
            case "province":
                getLocation(locationId, $scope.aProvince.cities, "city");
                break;
            case "city":
                getLocation(locationId, $scope.aCity.regions, "region");
                break;
            case "region":
                break;
        }
    };

    $modal.initBap = function () {

        if ( $('#modal-map-container').size() === 0) {
            console.warn("There IS NO #modal-map-container!");
            return;
        }
        $modal.bMap = new BMap.Map("modal-map-container"); // 创建地图实例
        $modal.bMap.addControl(new BMap.NavigationControl());
        $modal.bMap.addControl(new BMap.ScaleControl());
        $modal.bMap.addControl(new BMap.OverviewMapControl());

        var longitude = 0;
        var latitude = 0;

        var geolocation = new BMap.Geolocation();
        geolocation.getCurrentPosition(
            //获取位置信息成功
            function(position){
                if(this.getStatus() == BMAP_STATUS_SUCCESS){
                    console.warn('您的位置：' + position.point.lng + ',' + position.point.lat);
                    longitude = position.point.lng;
                    latitude  = position.point.lat;
                    // 根据坐标得到地址描述
                } else {
                    console.error('无法获取定位信息');
                }
            },{
                enableHighAccuracy: true,
                // timeout: 5000,
                maximumAge: 30*1000
            });

        var point = new BMap.Point(longitude, latitude);

        if ($scope.longitude !== 0 && $scope.latitude !== 0) {
            longitude = $scope.longitude;
            latitude = $scope.latitude;

            point = new BMap.Point(longitude, latitude);
            $scope.marker = new BMap.Marker(point);
            $scope.marker.setTitle($scope.name);
            $modal.bMap.addOverlay($scope.marker);
        }

        $modal.bMap.centerAndZoom(point, 12);  // 初始化地图，设置中心点坐标和地图级别

        // 创建地址解析器实例
        $modal.myGeo = new BMap.Geocoder();
        /**
         * 监听地图点击事件，获取点击处建筑物名称
         */
        $modal.bMap.addEventListener("click", function (e) {
            var pt = e.point;
            $scope.longitude = pt.lng;
            $scope.latitude = pt.lat;

            $modal.bMap.removeOverlay($scope.marker);
            $scope.marker.setPosition(new BMap.Point($scope.longitude, $scope.latitude));
            $modal.bMap.addOverlay($scope.marker);

            $modal.myGeo.getLocation(pt, function (rs) {
                var addComp = rs.addressComponents;
                /**
                 * 将获取到的建筑名赋值给$scope.address
                 */
                needSearch = false;

                var provinces = $filter('filter')($rootScope.locations, function (local) {
                    return (local.Name.substring(0, 2) === addComp.province.substring(0, 2));
                });

                if (provinces.length > 0) {
                    needLocation = false;
                    $scope.aProvince = provinces[0];
                    $scope.location = provinces[0];

                    var cities = $filter('filter')($scope.aProvince.cities, function (local) {
                        return (local.Name.substring(0, 2) === addComp.city.substring(0, 2));
                    });

                    if (cities.length > 0) {
                        needLocation = false;
                        $scope.aCity = cities[0];
                        $scope.location = cities[0];

                        var regions = $filter('filter')($scope.aCity.regions, function (local) {
                            return (local.Name.substring(0, 2) === addComp.district.substring(0, 2));
                        });

                        if (regions.length > 0) {
                            needLocation = false;
                            $scope.aRegion = regions[0];
                            $scope.location = regions[0];
                        }
                    }
                }

                //addComp.province + addComp.city + addComp.district
                $scope.address = addComp.street + addComp.streetNumber;
                /**
                 * 通知angularjs更新视图
                 */
                console.error("clicked addr:", $scope.address, $scope.longitude, $scope.latitude);
                $scope.$digest();
            });
        });
        /**
         * 初始化查询配置
         * @type {BMap.LocalSearch}
         */
        $modal.local = new BMap.LocalSearch($modal.bMap, {
            renderOptions: {
                map: $modal.bMap,
                panel: "results",
                autoViewport: true,
                selectFirstResult: true
            },
            pageCapacity: 1
        });
    };

    /**
     * 监听地址改变事件，当地址输入框的值改变时
     */
    $scope.$watch('address', function () {
        if (!needSearch) {
            console.error("!hasInit");
            needSearch = true;
            return
        }
        /**
         * 查询输入的地址并显示在地图上、调整地图视野
         */
        $modal.local.search($scope.location.LocationName + $scope.address);
        /**
         * 将输入的地址解析为经纬度
         */
        $modal.myGeo.getPoint($scope.address, function (point) {
            if (point) {
                /**
                 * 将地址解析为经纬度并赋值给$scope.longitude和$scope.latitude
                 */
                $scope.longitude = point.lng;
                $scope.latitude = point.lat;
            }
        });
    });

    $scope.$watch('location', function () {
        if (!needLocation ||
            !$scope.location || !$scope.location.LocationId === 0) {
            needLocation = true;
            return;
        }

        $modal.local.search($scope.location.LocationName);
        var res = $modal.local.getResults();
        console.warn("$modal.local.search($scope.location.LocationName);", res);
        var point = res.getPoi(0).point;
        if (point) {
            $scope.longitude = point.lng;
            $scope.latitude = point.lat;
        }
    });

    $uibModalInstance.opened.then(function () {
        $modal.initData();
    });

    $uibModalInstance.rendered.then(function () {
        $modal.initBap();
    });

    $modal.save = function () {
        console.info("ready to update bInfo!");
        $http.post("/boiler_update/?scope=location", {
            uid: $scope.uid,
            address: $scope.address,
            location_id: $scope.location.LocationId,
            longitude: $scope.longitude,
            latitude: $scope.latitude
        }).then(function (res) {
                console.error("Update bInfo Resp:", res);
                $rootScope.getBoilerList();
                swal({
                    title: "设备地址信息更新成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close('success');
                    currentData = null;
                });
            }, function (err) {
                swal({
                    title: "设备地址信息更新失败",
                    text: err.data,
                    type: "error"
                });
            });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        //currentData = null;
    };
});

angular.module('BoilerAdmin').controller('ModalBoilerInfoMaintainCtrl', function ($uibModalInstance,  $rootScope, $scope, $filter, $http, $log, currentData) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.data = {};
    $modal.data.inspectDate = {};

    var today = new Date();
    var nextMonth = today.setMonth(today.getMonth() + 1);
    console.error(nextMonth);
    var selectors = [
        'inner',
        'outer',
        'valve',
        'gauge'
    ];
    $scope.popup = {};

    $modal.initData = function () {
        if (currentData) {
            $modal.data.uid = currentData.Uid;
            $modal.data.name = currentData.Name;
            $modal.data.contact = currentData.Contact ? currentData.Contact.Name : "";
            $modal.data.phoneNumber = currentData.Contact ? currentData.Contact.PhoneNumber : "";
            $modal.data.mobileNumber = currentData.Contact ? currentData.Contact.MobileNumber : "";
            $modal.data.email = currentData.Contact ? currentData.Contact.Email : "";

            $modal.data.createdDate = moment(currentData.CreatedDate).format('YYYY-MM-DD');
            $modal.data.inspectDate.inner = currentData.InspectInnerDateNext ? new Date(currentData.InspectInnerDateNext) : nextMonth;
            $modal.data.inspectDate.outer = currentData.InspectOuterDateNext ? new Date(currentData.InspectOuterDateNext) : nextMonth;
            $modal.data.inspectDate.valve = currentData.InspectValveDateNext ? new Date(currentData.InspectValveDateNext) : nextMonth;
            $modal.data.inspectDate.gauge = currentData.InspectGaugeDateNext ? new Date(currentData.InspectGaugeDateNext) : nextMonth;
        } else {
            $modal.data.uid = "";
            $modal.data.name = "";
            $modal.data.contact = "";
            $modal.data.phoneNumber = "";
            $modal.data.mobileNumber = "";
            $modal.data.email = "";

            $modal.data.createdDate = moment(new Date()).format('YYYY-MM-DD');
            $modal.data.inspectDate.inner = nextMonth;
            $modal.data.inspectDate.outer = nextMonth;
            $modal.data.inspectDate.valve = nextMonth;
            $modal.data.inspectDate.gauge = nextMonth;
        }

        for (var i = 0; i < selectors.length; i++) {
            var selector = selectors[i];
            $scope.popup[selector] = {
                opened: false,
                minDate: today,
                initDate: $modal.data.inspectDate[selector] < today ? today : $modal.data.inspectDate[selector]
            }
        }
    };

    $uibModalInstance.opened.then(function () {
        $modal.initData();
        // console.error("BoilerMaintain:", $modal.data, $scope.popup);
    });

    $uibModalInstance.rendered.then(function () {
        // $modal.initBap();
    });

    $modal.save = function () {
        console.info("ready to update bInfo!");
        $http.post("/boiler_update/?scope=maintain", $modal.data)
            .then(function (res) {
            console.error("Update bInfo Resp:", res);
            $rootScope.getBoilerList();
            swal({
                title: "设备维护信息更新成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
            });
        }, function (err) {
            swal({
                title: "设备维护信息更新失败",
                text: err.data,
                type: "error"
            });
        });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');
        //currentData = null;
    };

    $scope.mytime = new Date();
    $scope.mytime.setMinutes(0);

    $scope.hstep = 1;
    $scope.mstep = 30;

    $scope.ismeridian = true;
    $scope.toggleMode = function() {
        $scope.ismeridian = ! $scope.ismeridian;
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
        initDate: new Date(),
        showWeeks: true
    };

    $scope.dateOptions = {
        formatYear: 'yy',
        maxDate: new Date(2020, 5, 22),
        minDate: new Date(),
        initDate: new Date(),
        startingDay: 0,

        "current-text": "今天"
    };

    $scope.open = function(selector) {
        $scope.popup[selector].opened = true;
    };

    $scope.setDate = function(year, month, day) {
        $scope.dt = new Date(year, month, day);
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

angular.module('BoilerAdmin').controller('ModalBoilerBindCtrl', function ($uibModalInstance, $rootScope, $http, $log, currentData) {
    var $modalBind = this;
    //console.warn("CurrentData In Modal:", currentData);
    $modalBind.currentData = currentData;
    $modalBind.name = currentData.Name;
    $modalBind.terminal = null;

    $http.get('/terminal_list/?scope=boiler_bind&boiler=' + currentData.Uid)
        .then(function (res) {
            // $scope.parameters = data;
            var terminals = [];
            console.error("Get Bind Terminal List Resp:", res);
            boiler_loop:
            for (var i in res.data) {
                var d = res.data[i];
                if (!d.Boilers) {
                    d.Boilers = [];
                }

                if (d.Boilers.length >= 8) {
                    continue;
                }

                for (var j in d.Boilers) {
                    var b = d.Boilers[j];
                    if (b.Uid === currentData.Uid) {
                        console.error("b.Uid === currentData.Uid", d.TerminalCode);
                        continue boiler_loop;
                    }
                }

                d.code = d.TerminalCode.toString();
                if (d.code.length < 6) {
                    for (var l = d.code.length; l < 6; l++) {
                        d.code = "0" + d.code;
                    }
                }

                d.text = "#" + d.code + " " + d.Name + "(机组" + (d.Boilers.length + 1) + ")";

                terminals.push(d);
            }

            if (terminals.length === 0) {
                terminals.push({Uid: "", text: "没有满足条件的终端"});
            } else {
                terminals.unshift({Uid: "", text: "请选择"});
            }

            $modalBind.terminals = terminals;
        });
    
    $modalBind.terminalFilter = function (value, index, array) {
        console.warn("terminalFilter", value, index, array);
        if (value.length >= 6) {
            return true;
        }

        return false;
    };

    $modalBind.ok = function () {
        console.info("ready to bind boiler!");
        $http.post("/boiler_bind/", {
            boiler_id: currentData.Uid,
            terminal_id: $modalBind.terminal.Uid
        }).then(function (res) {
            console.info("Update boilerBind Resp:", res);
            $rootScope.getBoilerList();
            swal({
                title: "绑定设备成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
            });
        }, function (err) {
            swal({
                title: "绑定设备失败",
                text: err.data,
                type: "error"
            });
        });
    };

    $modalBind.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        //currentData = null;
    };
});

// Please note that the close and dismiss bindings are from $uibModalInstance.

angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: '/directives/modal/terminal_config.html',
    bindings: {
        resolve: '<',
        close: '&',
        dismiss: '&'
    },
    controller: function () {
        var $ctrl = this;

        $ctrl.$onInit = function () {
            $ctrl.currentData = $ctrl.resolve.currentData;
            // $ctrl.selected = {
            //     item: $ctrl.items[0]
            // };
        };

        $ctrl.ok = function () {
            // $ctrl.close({$value: $ctrl.selected.item});
        };

        $ctrl.cancel = function () {
            $ctrl.dismiss({$value: 'cancel'});
        };
    }
});