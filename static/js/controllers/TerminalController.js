angular.module('BoilerAdmin').controller('TerminalController', function($rootScope, $scope, $http, $timeout, $uibModal, $log, $document,$state,$stateParams, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    terminal = this;
    terminal.isDone = false;
    terminal.msgData = {};

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    terminal.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
        //.withOption('rowCallback', rowCallback);

    terminal.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6),
        DTColumnDefBuilder.newColumnDef(7),
        DTColumnDefBuilder.newColumnDef(8).notSortable()
    ];

    terminal.refreshDataTables = function (callback) {
        $http.get('/terminal_issued_list')
            .then(function (res) {
                // $scope.parameters = data;
                var datasource = res.data;
                var num = 0;
                console.info("Get Terminal List Resp:", res);
                angular.forEach(datasource, function (d, key) {
                    var t = d.Terminal;
                    t.num = ++num;
                    t.code = t.TerminalCode.toString();
                    t.online = t.IsOnline? "在线" : "离线";
                   /* if(t.Boilers){
                        $http.get('/boiler/state/is_burning/?boiler=' + t.Boilers[0].Uid)
                            .then(function (res) {
                                // console.error("Fetch Status Resp:", res.data, d);
                                t.isBurning = (res.data.value === "true");
                                // t.online = (t.IsOnline||t.isBurning) ? "在线" : "离线";
                            }, function (err) {
                                console.error('Fetch Status Err!', err);
                            });
                    }*/

                    if(d.TermUpdateTime==="0001-01-01T00:00:00Z"){
                        d.TermUpdateTime = null;
                    }
                    if(d.PlatUpdateTime==="0001-01-01T00:00:00Z"){
                        d.PlatUpdateTime = null;
                    }

                    if (t.code.length < 6) {
                        for (var l = t.code.length; l < 6; l++) {
                            t.code = "0" + t.code;
                        }
                    }
                    t.simNum = t.SimNumber.length > 0 ? t.SimNumber : " - ";
                    t.ip = t.LocalIp.length > 0 ? t.LocalIp : " - ";


                    if (currentData && currentData.Uid === t.Uid) {
                        currentData = t;
                    }
                });

                terminal.datasource = datasource;
                console.info("Get Terminal List Resp:", terminal.datasource);

                if (callback) {
                    callback();
                }
                setTimeout(function () {
                    App.stopPageLoading();
                }, 1500);
                terminal.isDone = true;
            });
    };

    $scope.$on('modal.closing', function(event, reason, closed) {
        console.log('modal.closing: ' + (closed ? 'close' : 'dismiss') + '(' + reason + ')');
        var message = "You are about to leave the edit view. Uncaught reason. Are you sure?";
        switch (reason){
            // clicked outside
            case "backdrop click":
                message = "Any changes will be lost, are you sure?";
                break;

            // cancel button
            case "cancel":
                message = "Any changes will be lost, are you sure?";
                break;

            // escape key
            case "escape key press":
                message = "Any changes will be lost, are you sure?";
                break;
        }
        if (!confirm(message)) {
            event.preventDefault();
        }
    });

    var chanNum = 12;

    terminal.temperNum = chanNum;
    terminal.analogNum = chanNum;
    terminal.switchNum = 3;
    terminal.calcNum = chanNum;

    terminal.temperCols = [];
    terminal.analogCols = [];
    terminal.switchCols = [];
    terminal.calcCols = [];

    for (var i = 1; i <= chanNum; i++) {
        terminal.temperCols.push('Temper' + i + "_channel");
        terminal.analogCols.push('Analog' + i + "_channel");
        terminal.calcCols.push('C' + i + "_calculate_parm");
    }

    for (var i = 1; i <= 3; i++) {
        terminal.switchCols.push('Switch_' + i + "_channel");
    }

    /**
     * Origin Messages
     */

    // terminal.initMsgData = function () {
    //     terminal.msgData = {};
    //     terminal.msgData.isEmpty = true;
    //     terminal.msgData.lastUpload = null;
    // };

    terminal.viewMesData = function (data) {
        $state.go("terminal.message",{terminal:data});
        terminal.msgData.code = data;
    };

    //消息调试
    terminal.getOriginMessages = function () {
        // terminal.msgData = {};
        // terminal.msgData.isEmpty = true;
        // terminal.msgData.lastUpload = null;
        // Ladda.create(document.getElementById('terminal_origin_messages')).start();
        $http.get('/terminal_origin_message_list/?dev=origin&terminal=' + terminal.msgData.code)
            .then(function (res) {
                console.info("Get Terminal List Resp:", res);
                var datasource = res.data;
                if (!datasource) {
                    return;
                }

                for (var i = 0; i < datasource.length; i++) {
                    var d = datasource[i];
                    d.num = i;
                    d.date = new Date(d.TS);
                    for (var ic = 0; ic < terminal.calcCols.length; ic++) {
                        var col = terminal.calcCols[ic];
                        if (d[col]) {
                            d[col] = parseInt(d[col]);
                        }
                    }
                }

                terminal.msgData.datasource = datasource;

                if (datasource.length > 0){
                    terminal.msgData.isEmpty = false;
                }

                Ladda.create(document.getElementById('terminal_origin_messages')).stop();
            }, function (e) {
                Ladda.create(document.getElementById('terminal_origin_messages')).stop();
            });
    };

    /**
     * Modals
     * @type {[*]}
     */
    terminal.items = ['item1', 'item2', 'item3'];

    terminal.animationsEnabled = true;

    terminal.new = function () {
        currentData = null;
        editing = true;
        terminal.open('lg');
    };

    terminal.setting = function (data) {
        currentData = data;
        editing = true;
        terminal.open('lg');
    };

    terminal.view = function (data) {
        currentData = data;
        editing = false;
        terminal.open('lg');
    };

    terminal.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: terminal.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_config.html',
            controller: 'ModalTerminalCtrl',
            controllerAs: '$modal',
            backdrop:"static",
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                // items: function () {
                //     return terminal.items;
                // }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    terminal.channel = function (data,temp) {
        currentData = data;
        currentData.template = temp.Template;
        editing = true;
        terminal.channelOpen('lg');
        // for (var i = 0; i < maintain.datasource.length; i++) {
        //     if (maintain.datasource[i].Uid == uid) {
        //         currentData = maintain.datasource[i];
        //         //$log.info("maintain.edit GET:", currentData);
        //         maintain.open('lg');
        //         break;
        //     }
        // }
    };

    terminal.channelOpen = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: terminal.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_channel_config.html',
            controller: 'ModalTerminalChannelCtrl',
            controllerAs: '$modal',
            backdrop:"static",
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                // items: function () {
                //     return terminal.items;
                // }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    //终端批量配置
    terminal.groupConfig = function (){

        var modalInstance = $uibModal.open({
            templateUrl: 'groupConfig.html',
            controller: 'ModalGroupConfigCtrl',
            size: "lg",
            windowClass: 'zindex',
           /* resolve: {
                items1: function () {
                    return items;
                }
            }*/
        });

        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    };

    terminal.toggleAnimation = function () {
        terminal.animationsEnabled = !terminal.animationsEnabled;
    };

    //配置状态
    terminal.statusView = function (code) {
        $state.go("terminal.status",{terminal:code});
    };



    //表格横向滚动事件
    $scope.tableScroll = function(){
        var cdTableWrapper = angular.element(".cd-table-wrapper");//表格
        var scLeft = angular.element(".ter-table-container");
        var scLeftArrow = angular.element(document.getElementsByClassName("cd-scroll-right"));//向右箭头
        var scLeftArrow2 = angular.element(document.getElementsByClassName("cd-scroll-left"));//向左箭头
        scLeft.on("scroll", function() {
            //remove color gradient when table has scrolled to the end
            var total_table_width = parseInt(cdTableWrapper.css('width').replace('px', '')),
                table_viewport = parseInt(scLeft.css('width').replace('px', ''));

            if(scLeft.scrollLeft() >= total_table_width - table_viewport) {
                scLeft.addClass('table-end');
                scLeftArrow.css("display", "none");
            } else {
                scLeft.removeClass('table-end');
                scLeftArrow.css("display", "block");
                // $scope.scrollLeftEnd = true;
            };
            if(scLeft.scrollLeft()>0){
                scLeftArrow2.css("display", "block");
            }else{
                scLeftArrow2.css("display", "none");
            }
        });

        $scope.scRight = function() {
            // var scLeft = angular.element(".ter-table-container");
            var column_width = 200,
                new_left_scroll = parseInt(scLeft.scrollLeft()) + column_width;
            scLeft.animate({
                scrollLeft: new_left_scroll
            }, 200);
            // console.log(scLeft);
        };
        $scope.scLeft = function() {
            // var scLeft = angular.element(".ter-table-container");
            var column_width = 200,
                new_left_scroll = parseInt(scLeft.scrollLeft()) - column_width;
            scLeft.animate({
                scrollLeft: new_left_scroll
            }, 200);
        };

    };



    //下发
    //功能码
    $http.get("/term_function_code_list").then(function (res) {
        $rootScope.fcode = res.data;
    });

    //高低字节
    $http.get("/term_byte_list").then(function (res) {
        $rootScope.hlCodes = res.data;
    });



});

var terminal;
var currentData;
var editing;

angular.module('BoilerAdmin').controller('ModalTerminalCtrl', function ($uibModalInstance, $uibModal, $http, $log) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.editing = editing;
    $modal.editingCode = true;
    $modal.upgradeValue = "升级配置";


    //bin文件选择
    $http.get("/bin_list").then(function (res) {
        $modal.bins = res.data;
        console.log(res.data);
    });


    //升级配置
    $modal.upgrade = function () {
        $http.post("/upgrade_configuration",
            {path:$modal.bin.BinPath,uid:$modal.currentData.Uid})
            .then(function (res) {
                swal({
                    title: '升级成功，是否现在重启?',
                    text: res.data,
                    type: 'success',
                    showCancelButton: true,
                    confirmButtonText: '重启',
                    cancelButtonText: "取消",
                    // showLoaderOnConfirm: true
                }).then(function(isConfirm) {
                    App.startPageLoading({message: '正在加载数据...'});
                    if (isConfirm) {
                        $modal.sendConfMessage2();
                    }
                });
                },function (err) {
                swal({
                    title: "升级未成功...",
                    text:  err.data,
                    type: "warning"
                });
                // $modal.upgradeValue = err.data;


            })

    }

    //下发按钮
    $modal.down = function () {
        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/issued_config",
            {uid:$modal.currentData.Uid, code:$modal.currentData.code})
            .then(function (res) {
                App.stopPageLoading();
                swal({
                    title: "信息已发送",
                    text: res.data,
                    type: "success"
                });
                terminal.refreshDataTables();
            },function (err) {
                App.stopPageLoading();
                swal({
                    title: "失败",
                    text: err.data,
                    type: "warning"
                });
            });
    }



    console.log("currentData:",$modal.currentData);

    $modal.initCurrent = function () {
        if (currentData) {
            $modal.editingCode = false;

            $modal.title = "编辑参数";

            $modal.name = currentData.Name;
            $modal.code = currentData.code;
            $modal.org = currentData.Organization;

            $modal.boilers = currentData.Boilers;
            $modal.simNum = currentData.SimNumber;
            $modal.ip = currentData.LocalIp;
            $modal.upload = currentData.UploadFlag;
            $modal.uploadPeriod = currentData.UploadPeriod;

            $modal.description = currentData.Description;
            $modal.sets = [];

            console.error("Boilers:", currentData, currentData.Boilers);
            if (!currentData.Boilers) {
                currentData.Boilers = [];
            }

            for (var i = 0; i < 8; i++) {
                $modal.sets.push({
                    num: i + 1,
                    Name: "未配置",
                    hasDev: false
                });
            }

            for (var i = 0; i < currentData.Boilers.length; i++) {
                var boiler = currentData.Boilers[i];
                boiler.num = boiler.TerminalSetId;
                boiler.hasDev = true;
                $modal.sets[boiler.num - 1] = boiler;
            }

            console.error("Sets:", $modal.sets);

            $modal.deviceCount = currentData.Boilers.length;
        }
    };

    $modal.initCurrent();

    /*
     MyUidObject

     TerminalCode		int64		    `orm:"index"`
     LocalIp			string		    `orm:"size(60);null"`
     RemoteIp		    string		    `orm:"size(60);null"`
     RemotePort		    int
     UploadFlag		    bool
     UploadPeriod		int64
     SimNumber		    string		    `orm:"size(20)"`

     InstalledBy		string		    `orm:"size(60);null"`
     InstalledDate		time.Time	    `orm:"type(datetime);null"`

     Organization		*Organization	`orm:"rel(fk);null;index"`
     Boilers			[]*Boiler	    `orm:"reverse(many)"`
     */
    
    $modal.sendConfMessage = function () {
        var data = {
            uid: currentData.Uid,
            code: $modal.code,
            upload_flag: $modal.upload,
            upload_period: $modal.uploadPeriod
        };

        $http.post("/terminal_config/", data)
            .then(function (res) {
            console.warn("Send Terminal Config Message Done:", res);
            swal({
                title: "信息已发送",
                text: res.data,
                type: "success"
            });
        }, function (err) {
            swal({
                title: "信息发送失败",
                text: err.data,
                type: "error"
            });
        });
    };


    $modal.sendConfMessage2 = function () {
        var data = {
            uid: currentData.Uid
        };
        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/terminal_restart", data)
            .then(function (res) {
                console.warn("Send Terminal Config Message Done:", res);
                App.stopPageLoading();
                swal({
                    title: "信息已发送",
                    text: res.data,
                    type: "success"
                });
            }, function (err) {
                App.stopPageLoading();
                swal({
                    title: "信息发送失败",
                    text: err.data,
                    type: "error"
                });
            });
    };


    $modal.bindSet = function (set) {
        if (set.hasDev) {
            swal({
                title: "解除终端#" + currentData.code + "\n与该设备的绑定?",
                text: "解除绑定后，将无法收到来自 " + set.Name + " 的运行时数据。",
                type: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: '解绑',
                cancelButtonText: '取消'
            }).then(function () {
                $http.post("/boiler_unbind/", {
                    boiler_id: set.Uid,
                    terminal_id: currentData.Uid
                }).then(function (res) {
                    swal({
                        title: "绑定已解除",
                        text: "该终端已不再维护 " + set.Name + "相关数据，如需重新接入，请通过终端绑定流程进行再次绑定。",
                        type: "success"
                    });
                    terminal.refreshDataTables($modal.initCurrent);
                }, function (err) {
                    swal({
                        title: "解除绑定失败",
                        text: err.data,
                        type: "error"
                    });
                });
            });
        } else {
            $modal.openBind(set.num);
        }
    };

    $modal.ok = function () {
        if (!$modal.code.length || $modal.code.length !== 6) {
            return;
        }
        Ladda.create(document.getElementById('boiler_ok')).start();
        var ter = {
            uid: "",
            code: "",
            org_id: $modal.org ? $modal.org.Uid : "",
            name: $modal.name,
            sim_number: $modal.simNum,
            ip: $modal.ip,
            upload_flag: $modal.upload,
            upload_period: $modal.uploadPeriod,

            description: $modal.description
        };

        if (currentData) {
            ter.uid = currentData.Uid;
        }

        if ($modal.editingCode) {
            ter.code = $modal.code;
        }

        $http.post("/terminal_update/", ter)
            .then(function (res) {
            swal({
                title: "终端配置更新成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
                terminal.refreshDataTables();
            });
        }, function (err) {
            swal({
                title: "终端配置更新失败",
                text: err.data,
                type: "error"
            });
        });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.resetCode = function () {
        $modal.editingCode = true;
    };

    $modal.reboot = function () {
        var uid = null;
        if (currentData) {
            uid = currentData.Uid;
        }

        if (!uid || uid.length <= 0) {
            swal({
                title: "重启失败",
                text: "未知终端，无法重启",
                type: "error"
            });
            return;
        }

        swal({
            title: '确认重启该终端?',
            //text: "You won't be able to revert this!",
            type: 'question',
            showCancelButton: true,
            confirmButtonColor: '#d33',
            cancelButtonColor: '#3085d6',
            confirmButtonText: '确定',
            cancelButtonText: '取消'
        }).then(function () {
            $http.post("/terminal_reset/", {
                uid: uid
            }).then(function (res) {
                swal({
                    title: "终端已重启",
                    type: "success"
                });
            }, function (err) {
                swal({
                    title: "终端重启失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

    $modal.delete = function () {
        if (!currentData.Boilers || currentData.Boilers.length === 0) {
            swal({
                title: "确认删除该终端？",
                text: "注意：删除后将无法恢复，无法接收来自此终端的所有设备信息。",
                type: "warning",
                showCancelButton: true,
                //confirmButtonClass: "btn-danger",
                confirmButtonColor: "#d33",
                cancelButtonText: "取消",
                confirmButtonText: "删除",
                closeOnConfirm: false
            }).then(function () {
                $http.post("/terminal_delete/", {
                    uid: currentData.Uid
                }).then(function (res) {
                    swal({
                        title: "终端删除成功",
                        type: "success"
                    }).then(function () {
                        $uibModalInstance.close('success');
                        currentData = null;
                        terminal.refreshDataTables();
                    });
                }, function (err) {
                    swal({
                        title: "删除终端失败",
                        text: err.data,
                        type: "error"
                    });
                });
            });
        } else {
            swal({
                title: "无法删除该终端",
                text: "尚有" + currentData.Boilers.length + "台锅炉设备与该终端绑定，如需删除该终端，请先解绑其所有设备。",
                type: "error"
            });
        }

    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };

    $modal.openBind = function (setId, size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-body ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: terminal.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_bind.html',
            controller: 'ModalTerminalBindCtrl',
            controllerAs: '$modalBind',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex_sub',
            resolve: {
                $modal: function () {
                    return $modal;
                },
                currentTerminal: function () {
                    return currentData;
                },
                setId: function () {
                    return setId;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };
});

angular.module('BoilerAdmin').controller('ModalTerminalChannelCtrl', function ($uibModalInstance, $uibModal, $rootScope, $scope, $http, $log) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.editing = editing;
    $modal.editingCode = true;

    $modal.category = 9;

    //--------------下发----------------
    //功能码
    $modal.fcode = $rootScope.fcode; //分类
    $modal.fcode2 = [
        {Id: 1, Name: "01", Value: 1},
        {Id: 2, Name: "02", Value: 2},
        {Id: 3, Name: "03", Value: 3},
        {Id: 99, Name: "None", Value: 99}
    ];
    $modal.fcode1 = [
        {Id: 3, Name: "03", Value: 3},
        {Id: 4, Name: "04", Value: 4}

    ];

    //高低字节
    $modal.hlCodes = $rootScope.hlCodes; //分类

    //通信接口地址
    $http.get("/correspond_type_list").then(function (res) {
        $modal.communiInterfaces = res.data;
    });

    //数据位
    $http.get("/date_bit_list").then(function (res) {
        $modal.dataBits = res.data;
    });

    //心跳包频率
    $http.get("/heartbeat_packet_list").then(function (res) {
        $modal.heartbeats = res.data;
    });

    //校验位
    $http.get("/parity_bit").then(function (res) {
        $modal.checkDigits = res.data;
    });

    //从机地址
    $http.get("/slave_address_list").then(function (res) {
        $modal.subAdrs = res.data;
    });

    //停止位
    $http.get("/stop_bit_list").then(function (res) {
        $modal.stopBits = res.data;
    });


    //波特率
    $http.get("/baud_rate_list").then(function (res) {
        $modal.BaudRates = res.data;
    });


    $modal.terminalPass = "123456";

    //通信参数
    $modal.initParam = function () {
        if(!$modal.communParams){
            $http.post("/issued_communication",{terminal_code:currentData.code}).then(function (res) {
                $modal.communParams = res.data;
                // console.log($modal.communParams);
                //通信接口地址
                $modal.communiInterface  = $modal.communParams.CorrespondType;

                //数据位
                $modal.dataBit  = $modal.communParams.DataBit;

                //心跳包频率
                $modal.heartbeat = $modal.communParams.HeartBeat;

                //校验位
                $modal.checkDigit  = $modal.communParams.CheckBit;

                //从机地址
                $modal.subAdr  = $modal.communParams.SubAddress;

                //停止位
                $modal.stopBit  = $modal.communParams.StopBit;

                //波特率
                $modal.BaudRate = $modal.communParams.BaudRate;

            })
        }

    };


    //终端快速设置
    $modal.quickSet = function (){
        var items = [
            {
                id: $modal.currentData.code,
                template: $modal.currentData.template
            }
        ];
        console.log($modal.currentData);
        var modalInstance = $uibModal.open({
            templateUrl: '/directives/modal/terminal_channel_quick_set.html',
            controller: 'ModalQuickSetCtrl',
            size: "lg",
            windowClass: 'zindex_sub',
            resolve: {
                items1: function () {
                    return items;
                }
            }
        });


        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    };






    $modal.priorities = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];

    App.startPageLoading({message: '正在加载数据...'});
    $http.post('/channel_config_matrix/', {
        terminal_code: currentData.code
    }).then(function (res) {
        console.error("post /channel_config_matrix/ resp:", res);
        $modal.chanMatrix = res.data;
        $modal.dataMatrix = clone($modal.chanMatrix);

        for (var i = 0; i < $modal.chanMatrix.length; i++) {
            for (var j = 0; j < $modal.chanMatrix[i].length; j++) {

                if($modal.chanMatrix[i][j].RuntimeParameterChannelConfig){
                    $modal.chanMatrix[i][j].IsDefault = $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.IsDefault;
                    $modal.chanMatrix[i][j].Name = $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.Name;
                    $modal.chanMatrix[i][j].Parameter =  $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.Parameter;
                    $modal.chanMatrix[i][j].Status = $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.Status;
                    $modal.chanMatrix[i][j].Ranges = $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.Ranges;
                    $modal.chanMatrix[i][j].SwitchStatus = $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.SwitchStatus;
                    $modal.chanMatrix[i][j].SequenceNumber = $modal.chanMatrix[i][j].RuntimeParameterChannelConfig.SequenceNumber;
                    $modal.chanMatrix[i][j].noStatus = false;



                }



                    if($modal.chanMatrix[i][j].AnalogueSwitch.Modbus===0){
                        $modal.chanMatrix[i][j].AnalogueSwitch.Modbus = null;
                    }
                    if($modal.chanMatrix[i][j].AnalogueSwitch.BitAddress===0){
                        $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress = null;
                    }





                if (!$modal.chanMatrix[i][j].RuntimeParameterChannelConfig) {
                    $modal.chanMatrix[i][j] = {
                        Name: "默认(未配置)",
                        noStatus:true
                    }

                }

                if (((i !== 0 && i !== 1)||j !== 2 ) &&  (!$modal.dataMatrix[i][j].RuntimeParameterChannelConfig || $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.IsDefault) ) {
                    $modal.dataMatrix[i][j] = null;
                    $modal.chanMatrix[i][j].noStatus=true;
                } else {
                    $modal.dataMatrix[i][j].oParamId = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Parameter.Id;
                    $modal.dataMatrix[i][j].IsDefault = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.IsDefault;
                    $modal.dataMatrix[i][j].Name = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Name;
                    $modal.dataMatrix[i][j].Parameter =  $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Parameter;
                    $modal.dataMatrix[i][j].Status = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Status;
                    $modal.dataMatrix[i][j].Ranges = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Ranges;
                    $modal.dataMatrix[i][j].SwitchStatus = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.SwitchStatus;
                    $modal.dataMatrix[i][j].SequenceNumber = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.SequenceNumber;
                }
            }
        }

        setTimeout(function () {
            App.stopPageLoading();
        }, 800);
    });

    $modal.categoryChanged = function (category) {
        $modal.category = category;
    };

    $scope.setChannelConfStat = function () {

    };

    $modal.analogParameters = [{Id: 0, Name: '默认配置'}];
    $modal.switchParameters = [{Id: 0, Name: '默认配置'}];
    $modal.calculateParameters = [{Id: 0, Name: '默认配置'}];
    $modal.rangeParameters = [{Id: 0, Name: '默认配置'}];

    for (var i in $rootScope.parameters) {
        var param = $rootScope.parameters[i];
        switch (param.Category.Id) {
            case 10:
                $modal.analogParameters.push(param);
                break;
            case 11:
                $modal.switchParameters.push(param);
                break;
            case 12:
                $modal.calculateParameters.push(param);
                break;
            case 13:
                $modal.rangeParameters.push(param);
                break;
        }
    }

    $modal.parameters = [
        $modal.analogParameters,
        $modal.analogParameters,
        $modal.switchParameters,
        $modal.calculateParameters,
        $modal.rangeParameters
    ];






    $modal.matrixChanged = function (outerIndex, innerIndex) {
        console.info("Data Matrix:", $modal.dataMatrix, "\n", outerIndex, ":", innerIndex);
        if ($modal.dataMatrix[outerIndex][innerIndex].Parameter.Id === 0) {
            $modal.dataMatrix[outerIndex][innerIndex].Parameter = null;
            $modal.dataMatrix[outerIndex][innerIndex].oParamId = 0;
            $modal.dataMatrix[outerIndex][innerIndex].Status = -1;
            $modal.dataMatrix[outerIndex][innerIndex].SwitchStatus = 0;
            $modal.dataMatrix[outerIndex][innerIndex].Ranges = null;

            $modal.chanMatrix[outerIndex][innerIndex].AnalogueSwitch = null;
            // $modal.chanMatrix[outerIndex][innerIndex].Switch = null;
            $modal.chanMatrix[outerIndex][innerIndex].noStatus=true;
            // if($modal.chanMatrix[outerIndex][innerIndex].IsDefault!==true){
            //     $modal.chanMatrix[outerIndex][innerIndex].Name="默认(未配置)"
            // }
        } else {
            if ($modal.dataMatrix[outerIndex][innerIndex].oParamId !== $modal.dataMatrix[outerIndex][innerIndex].Parameter.Id) {
                $modal.dataMatrix[outerIndex][innerIndex].Ranges = [];
                $modal.dataMatrix[outerIndex][innerIndex].oParamId = $modal.dataMatrix[outerIndex][innerIndex].Parameter.Id;
            }

            if (!$modal.dataMatrix[outerIndex][innerIndex].Status || $modal.dataMatrix[outerIndex][innerIndex].Status === -1) {
                $modal.dataMatrix[outerIndex][innerIndex].Status = 0;
            }

            if ($modal.dataMatrix[outerIndex][innerIndex].Parameter.Category.Id === 11 && (!$modal.dataMatrix[outerIndex][innerIndex].SwitchStatus || $modal.dataMatrix[outerIndex][innerIndex].SwitchStatus === 0)) {
                $modal.dataMatrix[outerIndex][innerIndex].SwitchStatus = 1;
            }

            $modal.chanMatrix[outerIndex][innerIndex].noStatus = false;
            // if($modal.chanMatrix[outerIndex][innerIndex].Name==="默认(未配置)"){
            //     $modal.chanMatrix[outerIndex][innerIndex].Name = $modal.dataMatrix[outerIndex][innerIndex].Parameter.Name;
            // }
        }

    };

    $scope.matrixReset = function () {
        for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                if((i===0||i===1)&&j===2){
                    $modal.chanMatrix[i][j].AnalogueSwitch = null;
                }else {
                    $modal.dataMatrix[i][j] = null;
                    $modal.chanMatrix[i][j].AnalogueSwitch = null;
                    $modal.chanMatrix[i][j].noStatus=true;
                }

            }
        }
    };

    $modal.initCurrent = function () {
        if (currentData) {
            $modal.editingCode = false;

            $modal.title = "编辑参数";

            $modal.name = currentData.Name;
            $modal.code = currentData.code;
            $modal.boilers = currentData.Boilers;

            $modal.description = currentData.Description;
        }
    };

    $modal.initCurrent();


    $scope.fCodeChange =function (fcode,i,j) {
        console.log(fcode);
        if(j>=2&&j<5){
            if(fcode.Id ===1||fcode.Id ===2){
                $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress = 1;
            }
            if(fcode.Id ===99){
                $modal.chanMatrix[i][j].AnalogueSwitch.Modbus = 0;
                $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress = 0;
            }
        }


    };

    $scope.setStatus = function(outerIndex, innerIndex, status, sn) {
        // console.warn("$scope.setStatus", outerIndex, innerIndex, status, sn);
        $modal.dataMatrix[outerIndex][innerIndex].Status = status;
        if (status === 1) {
            $modal.dataMatrix[outerIndex][innerIndex].SequenceNumber = sn;
        } else {
            $modal.dataMatrix[outerIndex][innerIndex].SequenceNumber = -1;
        }
    };

    $scope.setSwitchStatus = function(outerIndex, innerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, innerIndex, status);
        $modal.dataMatrix[outerIndex][innerIndex].SwitchStatus = status;
    };

    $modal.openRange = function (currentChannel, number, size, parentSelector) {
        console.log(currentChannel);
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-body ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: terminal.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_channel_config_range.html',
            controller: 'ModalTerminalChannelConfigRangeCtrl',
            controllerAs: '$modalRange',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex_sub',
            resolve: {
                $modal: function () {
                    return $modal;
                },
                currentChannel: function () {
                    currentChannel.ChannelNumber = number;
                    return currentChannel;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    $modal.ok = function () {
        console.warn("$modal channel update!");
        if (!$modal.code.length || $modal.code.length !== 6) {
            console.error("$modal.code error:", $modal.code);
            return;
        }
        // Ladda.create(document.getElementById('channel_ok')).start();
        console.log("data:",$modal.dataMatrix ,"chan:",$modal.chanMatrix );
        var configUpload = [];
        for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                if ($modal.dataMatrix[i][j] !== $modal.chanMatrix[i][j]) {
                    // if (!$modal.dataMatrix[i][j]  && (($modal.chanMatrix[i][j] && $modal.chanMatrix[i][j].IsDefault === true))) {
                    //     console.warn('!!NULL data:', $modal.dataMatrix[i][j], $modal.chanMatrix[i][j]);
                    //     continue;
                    // }
                    var chanParamId = $modal.chanMatrix[i][j] && $modal.chanMatrix[i][j].Parameter ? $modal.chanMatrix[i][j].Parameter.Id : 0;
                    var dataParamId = $modal.dataMatrix[i][j] && $modal.dataMatrix[i][j].Parameter ? $modal.dataMatrix[i][j].Parameter.Id : 0;
                    var chanStatus = $modal.chanMatrix[i][j] ? $modal.chanMatrix[i][j].Status : 0 ;
                    var dataStatus = $modal.dataMatrix[i][j] ? $modal.dataMatrix[i][j].Status : 0 ;
                    var dataSeqNo = $modal.dataMatrix[i][j] && dataStatus === 1 ? $modal.dataMatrix[i][j].SequenceNumber : -1;
                    var chanSwitch = $modal.chanMatrix[i][j] ? $modal.chanMatrix[i][j].SwitchStatus : 0 ;
                    var dataSwitch = $modal.dataMatrix[i][j] ? $modal.dataMatrix[i][j].SwitchStatus : 0 ;
                    var chanRanges, dataRanges = [];
                    if (j === 5) {
                        chanRanges = $modal.chanMatrix[i][j] ? $modal.chanMatrix[i][j].Ranges : [] ;
                        dataRanges = $modal.dataMatrix[i][j] ? $modal.dataMatrix[i][j].Ranges : [] ;
                    }

                    if(dataParamId===0){
                        continue;
                    }

                    //功能码
                    var fcodeName = 0;
                    //MODBUS
                    var modbus = 0;
                    //位地址
                    var bitAddress = 0;
                    //高低字节
                    var termByte = 0;
                    if(j===0 || j===1 || j===5){
                        fcodeName = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Function ? $modal.chanMatrix[i][j].AnalogueSwitch.Function.Id:0;
                        modbus = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Modbus ? $modal.chanMatrix[i][j].AnalogueSwitch.Modbus:0;
                        termByte = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Byte?$modal.chanMatrix[i][j].AnalogueSwitch.Byte.Id:0 ;
                    }

                    if(j>=2 && j<5){
                        fcodeName = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Function?$modal.chanMatrix[i][j].AnalogueSwitch.Function.Id:0;
                        modbus = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Modbus? $modal.chanMatrix[i][j].AnalogueSwitch.Modbus:0;
                        bitAddress = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress? $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress:0;
                    }

                    if (dataParamId !== chanParamId || dataStatus !== chanStatus || chanSwitch !== dataSwitch || chanRanges !== dataRanges) {
                        var chan = j + 1;
                        var num = i + 1;
                        if (j >= 2 && j < 5) {
                            chan = 3;
                            num = i + (j - 2) * 16 + 1;
                        } else if (j === 5) {
                            chan = j;
                        }

                        var configData = {
                            terminal_code: $modal.code,
                            parameter_id: dataParamId,
                            channel_type: chan,
                            channel_number: num,

                            status: dataStatus,
                            sequence_number: dataSeqNo,

                            switch_status: dataSwitch,

                            fcodeName:fcodeName,
                            modbus:parseInt(modbus),
                            termByte:parseInt(termByte),
                            bitAddress:parseInt(bitAddress)

                        };


                        //表单验证
                        if(configData.parameter_id){
                            if(j===0 || j===1 || j===5){
                                if(fcodeName===0||modbus===0||termByte===0){
                                    swal({
                                        title: "通道配置更新失败",
                                        text:"配置信息不全 ，参数不能为0 "+ i + j,
                                        type: "error"
                                    });
                                    App.stopPageLoading();
                                    return false;
                                }
                                if(fcodeName===3){
                                    if(modbus<=40000||modbus>=50000){
                                        swal({
                                            title: "MODBUS地址错误",
                                            text:"功能码为03，MODBUS地址范围40001-49999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                                if(fcodeName===4){
                                    if(modbus<=30000||modbus>=40000){
                                        swal({
                                            title: "MODBUS地址错误",
                                            text:"功能码为04，MODBUS地址范围30001-39999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                            }

                            if(j>=2 && j<5){
                                if(j===2 && i===1){
                                    continue;
                                }
                                if(fcodeName!==99 && (fcodeName===0||modbus===0||bitAddress===0)){
                                    swal({
                                        title: "通道配置更新失败",
                                        text:"配置信息不全"+ i + j,
                                        type: "error"
                                    });
                                    App.stopPageLoading();
                                    return false;
                                }
                                if(fcodeName===1){
                                    if(modbus<1||modbus>=10000){
                                        swal({
                                            title: "开关通道MODBUS地址错误",
                                            text:"功能码为01，MODBUS地址范围00001-09999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                    if(bitAddress!=1){
                                        swal({
                                            title: "位地址错误",
                                            text:"功能码为01，对应位地址为1" ,
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                                if(fcodeName===2){
                                    if(modbus<=10000||modbus>=20000){
                                        swal({
                                            title: "开关通道MODBUS地址错误",
                                            text:"功能码为02，MODBUS地址范围10001-19999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                    if(bitAddress!=1){
                                        swal({
                                            title: "位地址错误",
                                            text:"功能码为02，对应位地址为1",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                                if(fcodeName===3){
                                    if(modbus<=40000||modbus>=50000){
                                        swal({
                                            title: "开关通道MODBUS地址错，请修改",
                                            text:"功能码为03，MODBUS地址范围40001-49999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                    if(bitAddress<1||bitAddress>16){
                                        swal({
                                            title: "位地址错误",
                                            text: "功能码为03，对应位地址范围为1-16",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                            }

                        }


                        if (j === 5 && dataParamId > 0) {
                            configData.ranges = [];
                            if (dataRanges.length <= 0) {
                                console.warn("data:", dataParamId, dataStatus, dataRanges);
                                console.warn("chan:", chanParamId, chanStatus, chanRanges);
                                swal({
                                    title: "状态量通道配置错误",
                                    text: "已配置的状态量通道，需要完成其状态值的配置才可提交",
                                    type: "error"
                                });
                                Ladda.create(document.getElementById('channel_ok')).stop();
                                return;
                            }
                            for (var n in dataRanges) {
                                var r = dataRanges[n];
                                var rg = {};
                                rg.min = r.Min;
                                rg.max = r.Max;
                                rg.name = r.Name;
                                switch (typeof n) {
                                    case "number":
                                        rg.value = n;
                                        break;
                                    case "string":
                                        rg.value = parseInt(n, 10);
                                        break;
                                }

                                configData.ranges.push(rg);
                            }
                        }

                        configUpload.push(configData);
                    }
                }
            }
        }



        var cParam = {
            terminal_code:$modal.code,
            baudRate : $modal.BaudRate?$modal.BaudRate.Id:0,
            dataBit : $modal.dataBit?$modal.dataBit.Id:0,
            stopBit : $modal.stopBit?$modal.stopBit.Id:0,
            checkDigit : $modal.checkDigit?$modal.checkDigit.Id:0,
            communInterface : $modal.communiInterface?$modal.communiInterface.Id:0,
            slaveAddress : $modal.subAdr?$modal.subAdr.Id:0,
            heartbeat:$modal.heartbeat?$modal.heartbeat.Id:0
        };

        if(!cParam.baudRate||!cParam.dataBit||!cParam.stopBit||!cParam.checkDigit||!cParam.communInterface||!cParam.slaveAddress||!cParam.heartbeat){
            swal({
                title: "通道配置更新失败",
                text:"通信参数不能为空 ",
                type: "error"
            });
            App.stopPageLoading();
            return false;
        }


        console.warn("$modal channel update!", configUpload);

        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/channel_config_update/", {channel:configUpload,param:cParam})
            .then(function (res) {
                 App.stopPageLoading();
                terminal.refreshDataTables();
                swal({
                    title: "通道配置更新成功，是否立刻下发？",
                    type: "success",
                        showCancelButton: true,
                        confirmButtonText: "确定下发",
                        cancelButtonText: "取消",
                        showLoaderOnConfirm:true
                    }).then(function(){
                        App.startPageLoading({message: '正在加载数据...'});
                        $http.post("/issued_config",
                            {uid:$modal.currentData.Uid, code:$modal.currentData.code})
                            .then(function (res) {
                                App.stopPageLoading();
                                swal({
                                    title: "信息已发送",
                                    text: res.data,
                                    type: "success"
                                });
                                terminal.refreshDataTables();
                            },function (err) {
                                App.stopPageLoading();
                                swal({
                                    title: "失败",
                                    text: err.data,
                                    type: "warning"
                                });
                            });
                    $uibModalInstance.close('success');
                    })
                    .then(function () {
                    currentData = null;
                });
            }, function (err) {
                swal({
                    title: "通道配置更新失败",
                    text: err.data,
                    type: "error"
                });
                App.stopPageLoading();
            });
        // Ladda.create(document.getElementById('channel_ok')).stop();

    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };

    //模板另存为
    $modal.templateSave = function () {
        if (!$modal.code.length || $modal.code.length !== 6) {
            console.error("$modal.code error:", $modal.code);
            return;
        }
        var configUpload = [];
        for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {

                if ($modal.dataMatrix[i][j] !== $modal.chanMatrix[i][j]) {
                    var chanParamId = $modal.chanMatrix[i][j] && $modal.chanMatrix[i][j].Parameter ? $modal.chanMatrix[i][j].Parameter.Id : 0;
                    var dataParamId = $modal.dataMatrix[i][j] && $modal.dataMatrix[i][j].Parameter ? $modal.dataMatrix[i][j].Parameter.Id : 0;
                    var chanStatus = $modal.chanMatrix[i][j] ? $modal.chanMatrix[i][j].Status : 0 ;
                    var dataStatus = $modal.dataMatrix[i][j] ? $modal.dataMatrix[i][j].Status : 0 ;
                    var dataSeqNo = $modal.dataMatrix[i][j] && dataStatus === 1 ? $modal.dataMatrix[i][j].SequenceNumber : -1;
                    var chanSwitch = $modal.chanMatrix[i][j] ? $modal.chanMatrix[i][j].SwitchStatus : 0 ;
                    var dataSwitch = $modal.dataMatrix[i][j] ? $modal.dataMatrix[i][j].SwitchStatus : 0 ;
                    var chanRanges, dataRanges = [];

                    if(dataParamId===0){
                        continue;
                    }


                    if (j === 5) {
                        chanRanges = $modal.chanMatrix[i][j] ? $modal.chanMatrix[i][j].Ranges : [] ;
                        dataRanges = $modal.dataMatrix[i][j] ? $modal.dataMatrix[i][j].Ranges : [] ;
                    }
                    //功能码
                    var fcodeName = 0;
                    //MODBUS
                    var modbus = 0;
                    //位地址
                    var bitAddress = 0;
                    //高低字节
                    var termByte = 0;
                    if(j===0 || j===1 || j===5){
                        fcodeName = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Function ? $modal.chanMatrix[i][j].AnalogueSwitch.Function.Id:0;
                        modbus = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Modbus ? $modal.chanMatrix[i][j].AnalogueSwitch.Modbus:0;
                        termByte = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Byte?$modal.chanMatrix[i][j].AnalogueSwitch.Byte.Id:0 ;
                    }

                    if(j>=2 && j<5){
                        fcodeName = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Function?$modal.chanMatrix[i][j].AnalogueSwitch.Function.Id:0;
                        modbus = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.Modbus? $modal.chanMatrix[i][j].AnalogueSwitch.Modbus:0;
                        bitAddress = $modal.chanMatrix[i][j].AnalogueSwitch && $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress? $modal.chanMatrix[i][j].AnalogueSwitch.BitAddress:0;
                    }

                    if (dataParamId !== chanParamId || dataStatus !== chanStatus || chanSwitch !== dataSwitch || chanRanges !== dataRanges) {
                        var chan = j + 1;
                        var num = i + 1;
                        if (j >= 2 && j < 5) {
                            chan = 3;
                            num = i + (j - 2) * 16 + 1;
                        } else if (j === 5) {
                            chan = j;
                        }

                        var configData = {
                            // terminal_code: $modal.code,
                            parameter_id: dataParamId,
                            channel_type: chan,
                            channel_number: num,

                            status: dataStatus,
                            sequence_number: dataSeqNo,

                            switch_status: dataSwitch,

                            fcodeName:fcodeName,
                            modbus:parseInt(modbus),
                            termByte:parseInt(termByte),
                            bitAddress:parseInt(bitAddress)

                        };

                        //表单验证
                        if(configData.parameter_id){
                            if(j===0 || j===1 || j===5){
                                if(fcodeName===0||modbus===0||termByte===0){
                                    swal({
                                        title: "通道配置更新失败",
                                        text:"配置信息不全 ，参数不能为0 "+ i + j,
                                        type: "error"
                                    });
                                    App.stopPageLoading();
                                    return false;
                                }
                                if(fcodeName===3){
                                    if(modbus<=40000||modbus>=50000){
                                        swal({
                                            title: "MODBUS地址错误",
                                            text:"功能码为03，MODBUS地址范围40001-49999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                                if(fcodeName===4){
                                    if(modbus<=30000||modbus>=40000){
                                        swal({
                                            title: "MODBUS地址错误",
                                            text:"功能码为04，MODBUS地址范围30001-39999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                            }

                            if(j>=2 && j<5){
                                if(j===2 && i===1){
                                    continue;
                                }
                                if(fcodeName!==99 && (fcodeName===0||modbus===0||bitAddress===0)){
                                    swal({
                                        title: "通道配置更新失败",
                                        text:"配置信息不全"+ i + j,
                                        type: "error"
                                    });
                                    App.stopPageLoading();
                                    return false;
                                }


                                if(fcodeName===1){
                                    if(modbus<1||modbus>=10000){
                                        swal({
                                            title: "开关通道MODBUS地址错误",
                                            text:"功能码为01，MODBUS地址范围00001-09999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                    if(bitAddress!=1){
                                        swal({
                                            title: "位地址错误",
                                            text:"功能码为01，对应位地址为1" ,
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                                if(fcodeName===2){
                                    if(modbus<=10000||modbus>=20000){
                                        swal({
                                            title: "开关通道MODBUS地址错误",
                                            text:"功能码为02，MODBUS地址范围10001-19999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                    if(bitAddress!=1){
                                        swal({
                                            title: "位地址错误",
                                            text:"功能码为02，对应位地址为1",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                                if(fcodeName===3){
                                    if(modbus<=40000||modbus>=50000){
                                        swal({
                                            title: "开关通道MODBUS地址错，请修改",
                                            text:"功能码为03，MODBUS地址范围40001-49999",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                    if(bitAddress<1||bitAddress>16){
                                        swal({
                                            title: "位地址错误",
                                            text: "功能码为03，对应位地址范围为1-16",
                                            type: "error"
                                        });
                                        App.stopPageLoading();
                                        return false;
                                    }
                                }
                            }

                        }

                        if (j === 5 && dataParamId > 0) {
                            configData.ranges = [];
                            if (dataRanges.length <= 0) {
                                console.warn("data:", dataParamId, dataStatus, dataRanges);
                                console.warn("chan:", chanParamId, chanStatus, chanRanges);
                                swal({
                                    title: "状态量通道配置错误",
                                    text: "已配置的状态量通道，需要完成其状态值的配置才可提交",
                                    type: "error"
                                });
                                Ladda.create(document.getElementById('channel_ok')).stop();
                                return;
                            }
                            for (var n in dataRanges) {
                                var r = dataRanges[n];
                                var rg = {};
                                rg.min = r.Min;
                                rg.max = r.Max;
                                rg.name = r.Name;
                                switch (typeof n) {
                                    case "number":
                                        rg.value = n;
                                        break;
                                    case "string":
                                        rg.value = parseInt(n, 10);
                                        break;
                                }

                                configData.ranges.push(rg);
                            }
                        }

                        configUpload.push(configData);
                    }
                }
            }
        }

        var cParam = {
            baudRate : $modal.BaudRate?$modal.BaudRate.Id:0,
            dataBit : $modal.dataBit?$modal.dataBit.Id:0,
            stopBit : $modal.stopBit?$modal.stopBit.Id:0,
            checkDigit : $modal.checkDigit?$modal.checkDigit.Id:0,
            communInterface : $modal.communiInterface?$modal.communiInterface.Id:0,
            slaveAddress : $modal.subAdr?$modal.subAdr.Id:0,
            heartbeat:$modal.heartbeat?$modal.heartbeat.Id:0
        };

        if(!cParam.baudRate||!cParam.dataBit||!cParam.stopBit||!cParam.checkDigit||!cParam.communInterface||!cParam.slaveAddress||!cParam.heartbeat){
            swal({
                title: "通道配置更新失败",
                text:"通信参数不能为空 ",
                type: "error"
            });
            App.stopPageLoading();
            return false;
        }

        var modalInstance = $uibModal.open({
            templateUrl: '/directives/modal/terminal_channel_template.html',
            controller: 'ModalTerminalTemplateCtrl',
            size: " ",
            windowClass: 'zindex_sub',
            resolve: {
                cParam: function () {
                    return cParam;
                },
                configUpload: function () {
                    return configUpload;
                },
                org:function () {
                    return currentData.Organization;
                }
            }
        });


        modalInstance.result.then(function (selectedItem) {
            // $scope.selected = selectedItem;
        }, function () {

        });
    }


});

angular.module('BoilerAdmin').controller('ModalTerminalBindCtrl', function ($uibModalInstance, $rootScope, $http, $filter, $modal, currentTerminal, setId) {
    var $modalBind = this;
    $modalBind.terminal = currentTerminal;
    $modalBind.code = currentTerminal.code;
    $modalBind.name = currentTerminal.Name;
    $modalBind.boiler = null;

    $modalBind.getBoilers = function () {
        var boilers = [];
        for (var i = 0; i < $rootScope.boilers.length; i++) {
            var b = $rootScope.boilers[i];
            var isMatched = false;
            for (var j = 0; j < currentTerminal.Boilers.length; j++) {
                if (b.Uid === currentTerminal.Boilers[j].Uid) {
                    isMatched = true;
                    break;
                }
            }

            if (isMatched) {
                continue;
            }

            b.org = "";
            if (b.Enterprise) {
                b.org = b.Enterprise.Name;
            } else if (b.Factory) {
                b.org = b.Factory.Name;
            } else if (b.Installed) {
                b.org = b.Installed.Name;
            }

            boilers.push(b);
        }
        if (boilers.length === 0) {
            boilers.push({Uid: "", Name: "没有未绑定的锅炉"});
        } else {
            boilers.unshift({Uid: "", Name: "请选择"});
        }
        $modalBind.boilers = boilers;
    };

    $modalBind.getBoilers();

    $rootScope.$watch('boilers', function () {
        $modalBind.getBoilers();
    });

    $modalBind.ok = function () {
        console.info("ready to bind boiler!", $modalBind.boiler, $modalBind.currentData);
        $http.post("/boiler_bind/", {
            boiler_id: $modalBind.boiler.Uid,
            terminal_id: $modalBind.terminal.Uid,
            terminal_set_id : setId
        }).then(function (res) {
            console.info("Update terminalBind Resp:", res);
            terminal.refreshDataTables($modal.initCurrent);
            swal({
                title: "绑定设备成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
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
    };
});

angular.module('BoilerAdmin').controller('ModalTerminalChannelConfigRangeCtrl', function ($uibModalInstance, $rootScope, $http, $filter, $modal, currentChannel) {
    var $modalRange = this;
    $modalRange.editing = editing;

    $modalRange.channel = currentChannel;
    $modalRange.number = currentChannel.ChannelNumber;

    $modalRange.ranges = clone(currentChannel.Ranges);
    if (!$modalRange.ranges) {
        $modalRange.ranges = [];
    }

    $modalRange.isValid = false;
    $modalRange.comment = "请完善相关信息";

    $modalRange.addNewRange = function () {
        $modalRange.ranges.push({});
    };

    $modalRange.removeRange = function (rg) {
        for (var i in $modalRange.ranges) {
            if (rg === $modalRange.ranges[i]) {
                $modalRange.ranges.splice(i, 1);
            }
        }
    };

    $modalRange.rangeChanged = function () {
        for (var i in $modalRange.ranges) {
            var rg = $modalRange.ranges[i];
            if (!rg.Min && typeof rg.Min !== "number" || rg.Min < 0 || rg.Min > 65535) {
                $modalRange.isValid = false;
                $modalRange.comment = "状态的范围低值不可为空，范围是 0-65535。";
                return;
            }

            if (!rg.Max && typeof rg.Max !== "number" || rg.Max < 0 || rg.Max > 65535) {
                $modalRange.isValid = false;
                $modalRange.comment = "状态的范围高值不可为空，范围是 0-65535。";
                return;
            }

            if (rg.Min > rg.Max) {
                $modalRange.isValid = false;
                $modalRange.comment = "状态的范围高值需大于或等于范围低值。";
                return;
            }

            if (i > 0 && rg.Min <= $modalRange.ranges[i - 1].Max) {
                $modalRange.isValid = false;
                $modalRange.comment = "状态间不可有值的交叉，后一个状态的低值不可小于或等于前一个状态的高值。";
                return;
            }

            if (!rg.Name || rg.Name.length <= 0) {
                $modalRange.isValid = false;
                $modalRange.comment = "状态的名称不可为空。";
                return;
            }

            $modalRange.isValid = true;
            $modalRange.comment = "配置正确";
        }
    };

    if ($modalRange.ranges.length <= 0) {
        $modalRange.addNewRange();
    }

    $modalRange.ok = function () {
        // swal({
        //     title: "修改状态量通道 #" + $modalRange.number + " 相关配置?",
        //     text: "修改后，将会覆盖之前该通道的状态量配置信息。",
        //     type: 'warning',
        //     showCancelButton: true,
        //     confirmButtonColor: '#d33',
        //     cancelButtonColor: '#3085d6',
        //     confirmButtonText: '修改',
        //     cancelButtonText: '取消'
        // }).then(function () {
        //
        // });
        $modalRange.channel.Uid = "";
        $modalRange.channel.Ranges = [];
        for (var i in $modalRange.ranges) {
            var rg = $modalRange.ranges[i];
            $modalRange.channel.Ranges.push(rg);
        }

        $uibModalInstance.dismiss('success');
    };

    $modalRange.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };

    $modalRange.rangeChanged();
});

//批量导入配置
angular.module('BoilerAdmin').controller('ModalGroupConfigCtrl', function ($scope, $uibModalInstance,$http) {
    $scope.items = [
        {
            start:null,
            end:null,
            template:null
        }
    ];
    $http.get("/template_list").then(function (res) {
        $scope.templates = res.data;
        console.log($scope.templates);
    });



    // $scope.selectedTemplate = $scope.template[0];
    $scope.addGroupConfig = function (){
        $scope.items.push({
            start:null,
            end:null,
            template:null});
    };

    $scope.removeGroupConfig = function (index) {
        $scope.items.splice(index,1);
    };

    $scope.ok = function () {
        if($scope.items.length<=0){
            swal({
                title: "没有配置数据",
                // text: err.data,
                type: "warning"
            });
            return false;
        }

        for(var i =0; i<$scope.items.length; i++){
            var template = $scope.items[i].template.Uid;
            $scope.items[i].template = template;
            // if($scope.items[i].start.length!==6||){
            //
            // }
        }



        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/template_group_config",{groupConfig:$scope.items}).then(function (res) {
            App.stopPageLoading();
            swal({
                title: "批量配置成功",
                text: res.data,
                type: "success"
            });
            terminal.refreshDataTables();
            $uibModalInstance.close();
        },function (err) {
            App.stopPageLoading();
            swal({
                title: "批量配置失败",
                text: err.data,
                type: "warning"
            });
        });

    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});

//快速设置
angular.module('BoilerAdmin').controller('ModalQuickSetCtrl', function ($scope, $uibModalInstance,$http, items1) {
    /*var items = [
        {
            id:$modal.currentData.code,
            template:"通用模板一"
        }
    ];*/
    $scope.items = items1;

    $http.get("/template_list").then(function (res) {
        $scope.templates = res.data;
        // console.log($scope.templates);
    });

    $scope.addQuickSet = function (){
        $scope.items.push({
            id: null,
            template:null});
    };

    $scope.removeQuickConfig = function (index) {
        $scope.items.splice(index,1);
    };

    $scope.ok = function () {
        $uibModalInstance.close($scope.selected.item);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});

//模板另存为
angular.module('BoilerAdmin').controller('ModalTerminalTemplateCtrl', function ($scope, $http, $uibModalInstance,cParam,configUpload,org) {

    console.log("cParam:",cParam,"configUpload:",configUpload,"org:",org);
    /*if(org==null){
        org={
            Uid:null
        }
    }*/
    $scope.templateName = "";
    $scope.ok = function () {
        if(org==null){
            swal({
                title: "终端没有所属企业",
                text: "请在 终端关联 中设置所属企业" ,
                type: "error"
            });
            return;
        }
        if(!$scope.templateName){
            swal({
                title: "模板未命名",
                text: "请填写需要保存的模板名称" ,
                type: "error"
            });
            return;
        }
        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/issued_template", {name:$scope.templateName,channel:configUpload,param:cParam,organizationUid:org.Uid})
            .then(function (res) {
                App.stopPageLoading();
                swal({
                    title: "模板保存成功",
                    text: res.data,
                    type: "success"
                }).then(function() {
                    $uibModalInstance.close('success');
                })
            }, function (err) {
                App.stopPageLoading();
                swal({
                    title: "模板保存失败",
                    text: err.data,
                    type: "error"
                });
            });
        // Ladda.create(document.getElementById('channel_ok')).stop();

    };


    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});


//配置状态
angular.module('BoilerAdmin').controller("terminalStatus",function ($scope,$http,$stateParams) {
    console.log($stateParams.terminal);
    $scope.statusList = [];
    $scope.refreshData = function () {
        $http.post("/terminal_error_list",{
            sn: $stateParams.terminal,
            startDate: $scope.startDate,
            endDate: $scope.endDate,
        }).then(function (res) {
            $scope.statusList = res.data;
            for(var i=0; i<$scope.statusList.length; i++){
                $scope.statusList[i].num = i+1;
                if($scope.statusList[i].ChannelNumber>0 && $scope.statusList[i].ChannelNumber<=12 ){
                    $scope.statusList[i].name = "模拟通道A 通道" + ($scope.statusList[i].ChannelNumber);
                }
                if($scope.statusList[i].ChannelNumber>12 && $scope.statusList[i].ChannelNumber<=24 ){
                    $scope.statusList[i].name = "模拟通道B 通道" + ($scope.statusList[i].ChannelNumber-12);
                }
                if($scope.statusList[i].ChannelNumber>24 && $scope.statusList[i].ChannelNumber<=40 ){
                    $scope.statusList[i].name = "开关通道1 通道" + ($scope.statusList[i].ChannelNumber-24);
                }
                if($scope.statusList[i].ChannelNumber>40 && $scope.statusList[i].ChannelNumber<=56 ){
                    $scope.statusList[i].name = "开关通道2 通道" + ($scope.statusList[i].ChannelNumber-40);
                }
                if($scope.statusList[i].ChannelNumber>56 && $scope.statusList[i].ChannelNumber<=72 ){
                    $scope.statusList[i].name = "开关通道3 通道" + ($scope.statusList[i].ChannelNumber-56);
                }
                if($scope.statusList[i].ChannelNumber>72 && $scope.statusList[i].ChannelNumber<=84 ){
                    $scope.statusList[i].name = "状态通道 通道" + ($scope.statusList[i].ChannelNumber-72);
                }
                if($scope.statusList[i].ChannelNumber===85){
                    $scope.statusList[i].name = "波特率";
                }
                if($scope.statusList[i].ChannelNumber===86){
                    $scope.statusList[i].name = "数据位";
                }
                if($scope.statusList[i].ChannelNumber===87){
                    $scope.statusList[i].name = "停止位";
                }
                if($scope.statusList[i].ChannelNumber===88){
                    $scope.statusList[i].name = "校验位";
                }
                if($scope.statusList[i].ChannelNumber===89){
                    $scope.statusList[i].name = "通讯接口类型";
                }
                if($scope.statusList[i].ChannelNumber===90){
                    $scope.statusList[i].name = "从机地址";
                }
                if($scope.statusList[i].ChannelNumber===91){
                    $scope.statusList[i].name = "心跳包频率";
                }
            }
            $scope.totalItems = $scope.statusList.length;
            $scope.currentPage = 1;
            $scope.pageNum = Math.ceil($scope.totalItems/20);
        },function (err) {

        })
    };




    //日期设置
    $scope.format = "yyyy-MM-dd";
    $scope.altInputFormats = ['yyyy-M!-d!'];
    $scope.startDate = new Date();
    $scope.endDate = new Date();



    $scope.popup1 = {
        opened: false
    };
    $scope.open1 = function () {
        $scope.popup1.opened = true;
    };

    $scope.popup2 = {
        opened: false
    };
    $scope.open2 = function () {
        $scope.popup2.opened = true;
    };



    $scope.setDataRange = function (range) {
        var startDate = new Date();
        var endDate = new Date();
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
        $scope.dataRange = range;
        $scope.startDate = startDate;
        $scope.endDate = endDate;

        $scope.dataRange = range;

        $scope.refreshData();
    };

    $scope.dateChanged = function () {
        console.warn("dateChanged:", $scope.startDate, "-", $scope.endDate);
        if ($scope.startDate < $scope.endDate) {
            $scope.refreshData();
        } else {
            $scope.statusList = [];
        }
    };


});





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
            // $ctrl.items = $ctrl.resolve.items;
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

function clone(obj) {
    var copy;

    // Handle the 3 simple types, and null or undefined
    if (null === obj || "object" !== typeof obj) return obj;

    // Handle Date
    if (obj instanceof Date) {
        copy = new Date();
        copy.setTime(obj.getTime());
        return copy;
    }

    // Handle Array
    if (obj instanceof Array) {
        copy = [];
        for (var i = 0, len = obj.length; i < len; i++) {
            copy[i] = clone(obj[i]);
        }
        return copy;
    }

    // Handle Object
    if (obj instanceof Object) {
        copy = {};
        for (var attr in obj) {
            if (obj.hasOwnProperty(attr)) copy[attr] = clone(obj[attr]);
        }
        return copy;
    }

    throw new Error("Unable to copy obj! Its type isn't supported.");
}