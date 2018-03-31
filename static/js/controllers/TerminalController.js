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
        $http.get('/terminal_list/')
            .then(function (res) {
                // $scope.parameters = data;
                var datasource = res.data;

                var num = 0;
                console.info("Get Terminal List Resp:", res);
                angular.forEach(datasource, function (d, key) {
                    d.num = ++num;
                    d.code = d.TerminalCode.toString();

                    d.online = d.IsOnline? "在线" : "离线";
                    if(d.Boilers){
                        $http.get('/boiler/state/is_burning/?boiler=' + d.Boilers[0].Uid)
                            .then(function (res) {
                                // console.error("Fetch Status Resp:", res.data, d);
                                d.isBurning = (res.data.value === "true");
                                d.online = (d.IsOnline||d.isBurning) ? "在线" : "离线";
                            }, function (err) {
                                console.error('Fetch Status Err!', err);
                            });
                    }


                    if (d.code.length < 6) {
                        for (var l = d.code.length; l < 6; l++) {
                            d.code = "0" + d.code;
                        }
                    }
                    d.simNum = d.SimNumber.length > 0 ? d.SimNumber : " - ";
                    d.ip = d.LocalIp.length > 0 ? d.LocalIp : " - ";


                    if (currentData && currentData.Uid === d.Uid) {
                        currentData = d;
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
        $state.go("terminal.message");
        terminal.msgData.code = data;
    }

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

    terminal.channel = function (data) {
        currentData = data;
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
        var items = [
            {
                start:680001,
                end:680100,
                template:"通用模板一"
            },
            {
                start:680001,
                end:680100,
                template:"通用模板一"
            }
        ];
        var modalInstance = $uibModal.open({
            templateUrl: 'groupConfig.html',
            controller: 'ModalGroupConfigCtrl',
            size: "lg",
            windowClass: 'zindex',
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
    }

    terminal.toggleAnimation = function () {
        terminal.animationsEnabled = !terminal.animationsEnabled;
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
                    title: "信息已发送",
                    text: res.data,
                    type: "success"
                });
                },function (err) {
                swal({
                    title: "升级未成功",
                    text: err.data,
                    type: "warning"
                });
            })
    }

    //按钮
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

        $http.post("/terminal_restart", data)
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

    //下发test


    //功能码
    $modal.fcode = $rootScope.fcode; //分类

    //高低字节
    $modal.hlCodes = $rootScope.hlCodes; //分类

    //通信接口地址
    $http.get("/correspond_type_list").then(function (res) {
        $modal.communiInterfaces = res.data;
    });
    $modal.communiInterface  = 0;

    //数据位
    $http.get("/date_bit_list").then(function (res) {
        $modal.dataBits = res.data;
    });
    $modal.dataBit  = 0;

    //心跳包频率
    $http.get("/heartbeat_packet_list").then(function (res) {
        $modal.heartbeats = res.data;
    });
    $modal.heartbeat = 0;

    //校验位
    $http.get("/parity_bit").then(function (res) {
        $modal.checkDigits = res.data;
    });
    $modal.checkDigit  = 0;

    //从机地址
    $http.get("/slave_address_list").then(function (res) {
        $modal.subAdrs = res.data;
    });
    $modal.subAdr  = 0;

    //停止位
    $http.get("/stop_bit_list").then(function (res) {
        $modal.stopBits = res.data;
    });
    $modal.stopBit  = 0;



    //波特率
    $http.get("/baud_rate_list").then(function (res) {
        $modal.BaudRates = res.data;
    });
    $modal.BaudRate  = 0;


    $modal.terminalPass = "123456";

    $modal.initParam = function () {
        // if(){
        //
        // }
        $http.post("/issued_communication",{code:currentData.code}).then(function (res) {
            console.log(res.data);
        })
    };


    //终端快速设置
    $modal.quickSet = function (){
        var items = [
            {
                id:680001,
                template:"通用模板一"
            },
            {
                id:680001,
                template:"通用模板一"
            }
        ];
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
    }






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


                   /*if(j===0||j===1||j===5){

                       $modal.fcodeName[i][j] = $modal.chanMatrix[i][j].Analogue.Function;
                       //MODBUS
                       $modal.mcode[i][j] = $modal.chanMatrix[i][j].Analogue.Modbus;
                       //高低字节
                       $modal.hlCodeNames[i][j] = $modal.chanMatrix[i][j].Analogue.Byte;

                   }
                   if(j>=2 && j<5){
                       $modal.fcodeName[i][j] = $modal.chanMatrix[i][j].Switch.Function;
                       //MODBUS
                       $modal.mcode[i][j] = $modal.chanMatrix[i][j].Switch.Modbus;
                       //位地址
                       $modal.bitAddress[i][j] = $modal.chanMatrix[i][j].Switch.BitAddress;
                   }*/


                }



                if (!$modal.chanMatrix[i][j].RuntimeParameterChannelConfig) {
                    $modal.chanMatrix[i][j] = {
                        Name: "默认(未配置)"
                    }

                }

                if (!$modal.dataMatrix[i][j].RuntimeParameterChannelConfig || $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.IsDefault) {
                    $modal.dataMatrix[i][j] = null;
                } else {
                    $modal.dataMatrix[i][j].oParamId = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Parameter.Id;
                    $modal.dataMatrix[i][j].IsDefault = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.IsDefault;
                    $modal.dataMatrix[i][j].Name = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Name;
                    $modal.dataMatrix[i][j].Parameter =  $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Parameter;
                    $modal.dataMatrix[i][j].Status = $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.Status;

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

            if($modal.chanMatrix[outerIndex][innerIndex].Name==="默认(未配置)"){
                $modal.chanMatrix[outerIndex][innerIndex].Name = $modal.dataMatrix[outerIndex][innerIndex].Parameter.Name;
            }
        }

    };

    $scope.matrixReset = function () {
        for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                $modal.dataMatrix[i][j] = null;
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
        Ladda.create(document.getElementById('channel_ok')).start();
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


                    //功能码
                    var fcodeName = 0;
                    //MODBUS
                    var modbus = 0;
                    //位地址
                    var bitAddress = 0;
                    //高低字节
                    var termByte = 0;
                    if(j===0 || j===1 || j===5){
                        fcodeName = $modal.chanMatrix[i][j].Analogue && $modal.chanMatrix[i][j].Analogue.Function ? $modal.chanMatrix[i][j].Analogue.Function.Id:0;
                        modbus = $modal.chanMatrix[i][j].Analogue&&$modal.chanMatrix[i][j].Analogue.Modbus ? $modal.chanMatrix[i][j].Analogue.Modbus:0;
                        termByte = $modal.chanMatrix[i][j].Analogue&&$modal.chanMatrix[i][j].Analogue.Byte?$modal.chanMatrix[i][j].Analogue.Byte.Id:0 ;
                    }

                    if(j>=2 && j<5){
                        fcodeName = $modal.chanMatrix[i][j].Switch && $modal.chanMatrix[i][j].Switch.Function?$modal.chanMatrix[i][j].Switch.Function.Id:0;
                        modbus = $modal.chanMatrix[i][j].Switch && $modal.chanMatrix[i][j].Switch.Modbus? $modal.chanMatrix[i][j].Switch.Modbus:0;
                        bitAddress = $modal.chanMatrix[i][j].Switch && $modal.chanMatrix[i][j].Switch.BitAddress? $modal.chanMatrix[i][j].Switch.BitAddress:0;
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



        console.warn("$modal channel update!", configUpload);

        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/channel_config_update/", {channel:configUpload,param:cParam})
            .then(function (res) {
                 App.stopPageLoading();

                swal({
                    title: "通道配置更新成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close('success');
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
        Ladda.create(document.getElementById('channel_ok')).stop();
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };
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


angular.module('BoilerAdmin').controller('ModalGroupConfigCtrl', function ($scope, $uibModalInstance, items1) {
    $scope.items = items1;

    $scope.template = [
        {
            id:1,
            name:"通用模板一"
        },
        {
            id:2,
            name:"通用模板二"
        },
        {
            id:3,
            name:"通用模板三"
        }
    ];
    $scope.selectedTemplate = $scope.template[0];
    $scope.addGroupConfig = function (){
        $scope.items.push({
            start:680001,
            end:680100,
            template:"通用模板一"});
    };


    $scope.ok = function () {
        $uibModalInstance.close($scope.selected.item);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});

angular.module('BoilerAdmin').controller('ModalQuickSetCtrl', function ($scope, $uibModalInstance, items1) {
    $scope.items = items1;

    $scope.template = [
        {
            id:1,
            name:"通用模板一"
        },
        {
            id:2,
            name:"通用模板二"
        },
        {
            id:3,
            name:"通用模板三"
        }
    ];
    $scope.selectedTemplate = $scope.template[0];
    $scope.addQuickSet = function (){
        $scope.items.push({
            id:680001,
            template:"通用模板一"});
    };


    $scope.ok = function () {
        $uibModalInstance.close($scope.selected.item);
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
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