angular.module('BoilerAdmin').controller("templateCtrl",function ($rootScope,$scope,$uibModal,$http,$filter, settings,DTOptionsBuilder, DTColumnDefBuilder) {
    template = this;

    App.initAjax();

    // set default layout mode
    $rootScope.settings.layout.pageContentWhite = true;
    $rootScope.settings.layout.pageBodySolid = true;
    $rootScope.settings.layout.pageSidebarClosed = false;

    template.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');

    template.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3).notSortable()
    ];

    //模板列表刷新
    template.refreshTemplate = function () {
        $http.get("/template_list").then(function (res) {
            var datasource = res.data;

            template.datasource = $filter("filter")(datasource,function (data) {
                if(!$rootScope.currentUser.Organization){
                    return true;
                }
               if(data.Organization.Uid === $rootScope.currentUser.Organization.Uid){
                   return true;
               }
            });
            for(var i = 0; i<template.datasource.length; i++){
                template.datasource[i].num = i+1;
            }
            // console.log(template.datasource);
        })
    };
    template.refreshTemplate();






    var currentData;
    var editing;

    template.new = function () {
        editing = true;
        var modalInstance = $uibModal.open({
            templateUrl: '/directives/modal/template_config.html',
            controller: 'ModalNewTemplateCtrl',
            controllerAs: '$modal',
            size: "lg",
            windowClass: 'zindex',
            resolve: {
                editing: function () {
                    return editing;
                }
            }
        });


        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    };

    template.edit = function (data) {
        currentData = data;
        editing = true;
        var modalInstance = $uibModal.open({
            templateUrl: '/directives/modal/template_config.html',
            controller: 'ModalEditTemplateCtrl',
            controllerAs: '$modal',
            size: "lg",
            windowClass: 'zindex',
            resolve: {
                currentData: function () {
                    return currentData;
                },
                editing: function () {
                    return editing;
                }
            }
        });


        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    }


    template.delete = function (uid) {
        swal({
            title: "确认删除该模板？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/template_delete", {
                uid: uid
            }).then(function (res) {
                swal({
                    title: "模板删除成功",
                    type: "success"
                }).then(function () {
                    template.refreshTemplate();
                });
            }, function (err) {
                swal({
                    title: "删除模板失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
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

var template;

angular.module('BoilerAdmin').controller('ModalEditTemplateCtrl', function ($rootScope,$scope,$uibModal, $uibModalInstance,$http,currentData,editing) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.editing = editing;
    $modal.editingCode = true;

    $modal.category = 9;

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

    $modal.priorities = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];


    //模拟通道1
    $modal.initAnalog1=function () {
        if(!$modal.analogData1){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_analog_one",{
                uid:currentData.Uid
            }).then(function (res) {
                $modal.analogData1 = res.data;
                for (var i=0; i<$modal.analogData1.length; i++){
                    if($modal.analogData1[i].Modbus===0){
                        $modal.analogData1[i].Modbus = null;
                    }
                }
            });
            setTimeout(function () {
                App.stopPageLoading();
            },800);

        }

    };

    //模拟通道2
    $modal.initAnalog2=function () {
        if(!$modal.analogData2){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_analog_two",{
                uid:currentData.Uid
            }).then(function (res) {
                $modal.analogData2 = res.data;
                for (var i=0; i<$modal.analogData2.length; i++){
                    if($modal.analogData2[i].Modbus===0){
                        $modal.analogData2[i].Modbus = null;
                    }
                }
            });
            setTimeout(function () {
                App.stopPageLoading();
            },800);
        }
    };

    //开关通道1
    $modal.initSwitch1=function () {
        if(!$modal.switchData1){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_switch_one",{
                uid:currentData.Uid
            }).then(function (res) {
                $modal.switchData1 = res.data;
                for(var i = 0; i<$modal.switchData1.length; i++){
                    if($modal.switchData1[i].Modbus===0){
                        $modal.switchData1[i].Modbus = null;
                    }
                    if($modal.switchData1[i].BitAddress===0){
                        $modal.switchData1[i].BitAddress = null;
                    }
                    if(i===1){
                        $modal.switchData1[i].Func={Id: 99, Name: "None", Value: 99}
                    }
                }

                // console.log($modal.switchData1);
            });
            setTimeout(function () {
                App.stopPageLoading();
            },800);
        }
    };

    //开关通道2
    $modal.initSwitch2=function () {
        if(!$modal.switchData2){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_switch_two",{
                uid:currentData.Uid
            }).then(function (res) {
                $modal.switchData2 = res.data;
                for(var i = 0; i<$modal.switchData2.length; i++){
                    if($modal.switchData2[i].Modbus===0){
                        $modal.switchData2[i].Modbus = null;
                    }
                    if($modal.switchData2[i].BitAddress===0){
                        $modal.switchData2[i].BitAddress = null;
                    }
                }
            });
            setTimeout(function () {
                App.stopPageLoading();
            },800);
        }
    };

    //开关通道3
    $modal.initSwitch3=function () {
        if(!$modal.switchData3){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_switch_Three",{
                uid:currentData.Uid
            }).then(function (res) {
                $modal.switchData3 = res.data;
                for(var i = 0; i<$modal.switchData3.length; i++){
                    if($modal.switchData3[i].Modbus===0){
                        $modal.switchData3[i].Modbus = null;
                    }
                    if($modal.switchData3[i].BitAddress===0){
                        $modal.switchData3[i].BitAddress = null;
                    }
                }
            });
            setTimeout(function () {
                App.stopPageLoading();
            },800);
        }
    };

    //状态通道
    $modal.initRange=function () {
        if(!$modal.rangData){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_range",{
                uid:currentData.Uid
            }).then(function (res) {
                $modal.rangData = res.data;
                for (var i=0; i<$modal.rangData.length; i++){
                    if($modal.rangData[i].Modbus===0){
                        $modal.rangData[i].Modbus = null;
                    }
                }
            });
            setTimeout(function () {
                App.stopPageLoading();
            },800);
        }
    };

    //通信参数
    $modal.initParam = function () {
        if(!$modal.communParams){
            App.startPageLoading({message: '正在加载数据...'});
            $http.post("/template_communication",{
                uid:currentData.Uid
            }).then(function (res) {
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
            setTimeout(function () {
                App.stopPageLoading();
            },800);
        }

    };




    $modal.categoryChanged = function (category) {
        $modal.category = category;
        if($modal.category === 16){
            $modal.initAnalog1();
            $modal.initAnalog2();
            $modal.initSwitch1();
            $modal.initSwitch2();
            $modal.initSwitch3();
            $modal.initRange();
        }
    };

    //运行参数列表导入
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

    $modal.dataChanged = function (data) {
        if (data.Parameter.Id === 0) {
            data.Parameter = null;
            data.oParamId = 0;
            data.Status = -1;
            data.SwitchStatus = 0;
            data.Ranges = null;
            data.Func = null;
            data.Byte = null;
            data.Modbus = null;
            data.BitAddress = null;
        } else {
           /* if ($modal.dataMatrix[outerIndex][innerIndex].oParamId !== $modal.dataMatrix[outerIndex][innerIndex].Parameter.Id) {
                $modal.dataMatrix[outerIndex][innerIndex].Ranges = [];
                $modal.dataMatrix[outerIndex][innerIndex].oParamId = $modal.dataMatrix[outerIndex][innerIndex].Parameter.Id;
            }*/

            if (!data.Status || data.Status === -1) {
                data.Status = 0;
            }

            if (data.Parameter.Category.Id === 11 && (!data.SwitchStatus || data.SwitchStatus === 0)) {
                data.SwitchStatus = 1;
            }

            console.log(data);
        }

    };

    //恢复默认
    $scope.matrixReset = function () {
        /*for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                $modal.dataMatrix[i][j] = null;
            }
        }*/

        var data1 = function (data) {
            for (var i=0; i<12; i++){
                data[i]={
                    BitAddress: 0,
                    Byte: null,
                    ChannelNumber: 0,
                    ChannelType: 0,
                    Func: null,
                    Modbus: null,
                    Parameter: null,
                    Ranges: null,
                    SequenceNumber: 0,
                    Status: 0,
                    SwitchStatus: 0,
                    Template: null
                };
            }
        };
        var data2 = function (data) {
            for (var i=0; i<16; i++){
                data[i]={
                    BitAddress: null,
                    Byte: null,
                    Func: null,
                    Modbus: null,
                    Parameter: null,
                    Ranges: null,
                    SequenceNumber: 0,
                    Status: 0,
                    SwitchStatus: 0
                };
            }
        }

        data1($modal.analogData1);
        data1($modal.analogData2);
        data2($modal.switchData1);
        data2($modal.switchData2);
        data2($modal.switchData3);
        data1($modal.rangData);

    };

    $modal.initCurrent = function () {
        if (currentData) {
            $modal.editingCode = false;


            $modal.name = currentData.Name;
            $modal.code = currentData.code;
            $modal.boilers = currentData.Boilers;

            $modal.description = currentData.Description;
        }
    };

    $modal.initCurrent();

    $scope.fCodeChange1 =function (fcode,i) {
        // console.log(fcode);
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.switchData1[i].BitAddress = 1;
        }
        if(fcode.Id ===99){
            $modal.switchData1[i].Modbus = 0;
            $modal.switchData1[i].BitAddress = 0;
        }
    };
    $scope.fCodeChange2 =function (fcode,i) {
        // console.log(fcode);
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.switchData2[i].BitAddress = 1;
        }
        if(fcode.Id ===99){
            $modal.switchData2[i].Modbus = 0;
            $modal.switchData2[i].BitAddress = 0;
        }
    };
    $scope.fCodeChange3 =function (fcode,i) {
        // console.log(fcode);
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.switchData3[i].BitAddress = 1;
        }
        if(fcode.Id ===99){
            $modal.switchData3[i].Modbus = 0;
            $modal.switchData3[i].BitAddress = 0;
        }
    };


    //位置设置
    $scope.setStatus1 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.analogData1[index].Status = status;
        if (status === 1) {
            $modal.analogData1[index].SequenceNumber = sn;
        } else {
            $modal.analogData1[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus2 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.analogData2[index].Status = status;
        if (status === 1) {
            $modal.analogData2[index].SequenceNumber = sn;
        } else {
            $modal.analogData2[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus3 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.switchData1[index].Status = status;
        if (status === 1) {
            $modal.switchData1[index].SequenceNumber = sn;
        } else {
            $modal.switchData1[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus4 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.switchData2[index].Status = status;
        if (status === 1) {
            $modal.switchData2[index].SequenceNumber = sn;
        } else {
            $modal.switchData2[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus5 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.switchData3[index].Status = status;
        if (status === 1) {
            $modal.switchData3[index].SequenceNumber = sn;
        } else {
            $modal.switchData3[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus6 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.rangData[index].Status = status;
        if (status === 1) {
            $modal.rangData[index].SequenceNumber = sn;
        } else {
            $modal.rangData[index].SequenceNumber = -1;
        }
    };


    //状态设置
    $scope.setSwitchStatus1 = function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $modal.switchData1[outerIndex].SwitchStatus = status;
    };
    $scope.setSwitchStatus2 = function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $modal.switchData2[outerIndex].SwitchStatus = status;
    };
    $scope.setSwitchStatus3 = function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $modal.switchData3[outerIndex].SwitchStatus = status;
    };

    $modal.openRange = function (currentChannel, number, size, parentSelector) {
        console.log(currentChannel,number);
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-body ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: template.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_channel_config_range.html',
            controller: 'ModalTemplateRangeCtrl',
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
                },
                editing: function () {
                    return $modal.editing;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            console.log("Modal dismissed");
        });
    };



    $modal.ok = function () {
        $modal.templatesChannel = [
            $modal.analogData1,
            $modal.analogData2,
            $modal.switchData1,
            $modal.switchData2,
            $modal.switchData3,
            $modal.rangData
        ];
        console.log($modal.templatesChannel);
        var configUpload = [];

        for(var i = 0; i < $modal.templatesChannel.length; i++){
            for(var j = 0; j < $modal.templatesChannel[i].length; j++){
                var paramId = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Parameter ? $modal.templatesChannel[i][j].Parameter.Id : 0;
                var status = $modal.templatesChannel[i][j] ? $modal.templatesChannel[i][j].Status : 0 ;
                var seqNo = $modal.templatesChannel[i][j] && status === 1 ? $modal.templatesChannel[i][j].SequenceNumber : -1;
                var dataSwitch = $modal.templatesChannel[i][j] ? $modal.templatesChannel[i][j].SwitchStatus : 0 ;
                var ranges = [];
                var chan = i + 1;
                var num = j + 1;
                //功能码
                var fcodeName = 0;
                //MODBUS
                var modbus = 0;
                //位地址
                var bitAddress = 0;
                //高低字节
                var termByte = 0;

                if(paramId===0){
                    continue;
                }

                if(i===0||i===1){
                    chan = i + 1;
                    num = j + 1;
                    fcodeName = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Func ? $modal.templatesChannel[i][j].Func.Id:0;
                    modbus = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Modbus ? $modal.templatesChannel[i][j].Modbus:0;
                    termByte = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Byte? $modal.templatesChannel[i][j].Byte.Id:0 ;
                }
                if(i>=2 && i <5){
                    chan = 3;
                    num = j + (i - 2) * 16 + 1;
                    fcodeName = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Func ? $modal.templatesChannel[i][j].Func.Id:0;
                    modbus = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Modbus? $modal.templatesChannel[i][j].Modbus:0;
                    bitAddress = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].BitAddress? $modal.templatesChannel[i][j].BitAddress:0;
                }
                if(i===5){
                    chan = 5;
                    ranges = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Ranges ? $modal.templatesChannel[i][j].Ranges : [];
                    fcodeName = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Func ? $modal.templatesChannel[i][j].Func.Id:0;
                    modbus = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Modbus ? $modal.templatesChannel[i][j].Modbus:0;
                    termByte = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Byte? $modal.templatesChannel[i][j].Byte.Id:0 ;;
                }




                var configData = {
                    parameter_id: paramId,
                    channel_type: chan,
                    channel_number: num,
                    status: status,
                    sequence_number: seqNo,
                    switch_status: dataSwitch,
                    fcodeName:fcodeName,
                    modbus:parseInt(modbus),
                    termByte:parseInt(termByte),
                    bitAddress:parseInt(bitAddress)
                };

                //表单验证
                if(configData.parameter_id){
                    if(i===0 || i===1 || i===5){
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

                    if(i>=2 && i<5){
                        if(i===2 && j===1){
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


                if (i === 5 && paramId > 0) {
                    configData.ranges = [];
                    if (ranges.length <= 0) {
                        console.warn("data:", paramId, status, ranges);
                        swal({
                            title: "状态量通道配置错误",
                            text: "已配置的状态量通道，需要完成其状态值的配置才可提交",
                            type: "error"
                        });
                        return;
                    }
                    for (var n in ranges) {
                        var r = ranges[n];
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


        console.warn("$modal channel update!", configUpload);

        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/template_update", {TemplateUpdate:{uid:currentData.Uid,chan:configUpload,param:cParam,name:$modal.currentData.Name}})
            .then(function (res) {
                App.stopPageLoading();
                swal({
                    title: "模板配置更新成功",
                    type: "success"
                }).then(function () {
                    template.refreshTemplate();
                    $uibModalInstance.close('success');
                    currentData = null;
                });
            }, function (err) {
                App.stopPageLoading();
                swal({
                    title: "模板配置更新失败",
                    text: err.data,
                    type: "error"
                });
            });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');
        currentData = null;
    };
});



angular.module('BoilerAdmin').controller('ModalNewTemplateCtrl', function ($rootScope,$scope, $uibModalInstance,$http,editing) {
    var $modal = this;
    $modal.currentData = null;
    $modal.editing = editing;
    $modal.editingCode = true;

    $modal.category = 9;

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

    $modal.priorities = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];


    //模拟通道1
    $modal.initAnalog1=function () {
        $modal.analogData1 = [ ];
        for (var i=0; i<12; i++){
            $modal.analogData1[i]={
                BitAddress: 0,
                Byte: null,
                ChannelNumber: 0,
                ChannelType: 0,
                Func: null,
                Modbus: null,
                Parameter: null,
                Ranges: null,
                SequenceNumber: 0,
                Status: 0,
                SwitchStatus: 0,
                Template: null
            };
        }

    };

    //模拟通道2
    $modal.initAnalog2=function () {
        $modal.analogData2 = [ ];
        for (var i=0; i<12; i++){
            $modal.analogData2[i]={
                BitAddress: 0,
                Byte: null,
                ChannelNumber: 0,
                ChannelType: 0,
                Func: null,
                Modbus: null,
                Parameter: null,
                Ranges: null,
                SequenceNumber: 0,
                Status: 0,
                SwitchStatus: 0,
                Template: null
            };
        }
    };

    //开关通道1
    $modal.initSwitch1=function () {
        $modal.switchData1 = [ ];
        for (var i=0; i<16; i++){
            $modal.switchData1[i]={
                BitAddress: null,
                Byte: null,
                Func: null,
                Modbus: null,
                Parameter: null,
                Ranges: null,
                SequenceNumber: 0,
                Status: 0,
                SwitchStatus: 0
            };
        }

    };

    //开关通道2
    $modal.initSwitch2=function () {
        $modal.switchData2 = [ ];
        for (var i=0; i<16; i++){
            $modal.switchData2[i]={
                BitAddress: null,
                Byte: null,
                Func: null,
                Modbus: null,
                Parameter: null,
                Ranges: null,
                SequenceNumber: 0,
                Status: 0,
                SwitchStatus: 0
            };
        }
    };

    //开关通道3
    $modal.initSwitch3=function () {
        $modal.switchData3 = [ ];
        for (var i=0; i<16; i++){
            $modal.switchData3[i]={
                BitAddress: null,
                Byte: null,
                Func: null,
                Modbus: null,
                Parameter: null,
                Ranges: null,
                SequenceNumber: 0,
                Status: 0,
                SwitchStatus: 0
            };
        }
    };

    //状态通道
    $modal.initRange=function () {
        $modal.rangData = [ ];
        for (var i=0; i<12; i++){
            $modal.rangData[i]={
                BitAddress: 0,
                Byte: null,
                ChannelNumber: 0,
                ChannelType: 0,
                Func: null,
                Modbus: null,
                Parameter: null,
                Ranges: null,
                SequenceNumber: 0,
                Status: 0,
                SwitchStatus: 0,
                Template: null
            };
        }


    };

    //通信参数
    $modal.initParam = function () {
        //通信接口地址
        $modal.communiInterface  = null;

        //数据位
        $modal.dataBit  = null;

        //心跳包频率
        $modal.heartbeat = null;

        //校验位
        $modal.checkDigit  = null;

        //从机地址
        $modal.subAdr  = null;

        //停止位
        $modal.stopBit  = null;

        //波特率
        $modal.BaudRate = null;

    };



    $modal.categoryChanged = function (category) {
        $modal.category = category;
    };

    //运行参数列表导入
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

    $modal.dataChanged = function (data) {
        if (data.Parameter.Id === 0) {
            data.Parameter = null;
            data.oParamId = 0;
            data.Status = -1;
            data.SwitchStatus = 0;
            data.Ranges = null;
        } else {
            /* if ($modal.dataMatrix[outerIndex][innerIndex].oParamId !== $modal.dataMatrix[outerIndex][innerIndex].Parameter.Id) {
                 $modal.dataMatrix[outerIndex][innerIndex].Ranges = [];
                 $modal.dataMatrix[outerIndex][innerIndex].oParamId = $modal.dataMatrix[outerIndex][innerIndex].Parameter.Id;
             }*/

            if (!data.Status || data.Status === -1) {
                data.Status = 0;
            }

            if (data.Parameter.Category.Id === 11 && (!data.SwitchStatus || data.SwitchStatus === 0)) {
                data.SwitchStatus = 1;
            }

            console.log(data);
        }

    };

    //恢复默认
    /* $scope.matrixReset = function () {
         for (var i = 0; i < $modal.dataMatrix.length; i++) {
             for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                 $modal.dataMatrix[i][j] = null;
                 if($modal.chanMatrix[i][j].IsDefault!=true){
                     $modal.chanMatrix[i][j].Name="默认(未配置)"
                 }
             }
         }
     };*/



    $scope.fCodeChange1 =function (fcode,i) {
        // console.log(fcode);
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.switchData1[i].BitAddress = 1;
        }
    };
    $scope.fCodeChange2 =function (fcode,i) {
        // console.log(fcode);
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.switchData2[i].BitAddress = 1;
        }
    };
    $scope.fCodeChange3 =function (fcode,i) {
        // console.log(fcode);
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.switchData3[i].BitAddress = 1;
        }
    };


    //位置设置
    $scope.setStatus1 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.analogData1[index].Status = status;
        if (status === 1) {
            $modal.analogData1[index].SequenceNumber = sn;
        } else {
            $modal.analogData1[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus2 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.analogData2[index].Status = status;
        if (status === 1) {
            $modal.analogData2[index].SequenceNumber = sn;
        } else {
            $modal.analogData2[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus3 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.switchData1[index].Status = status;
        if (status === 1) {
            $modal.switchData1[index].SequenceNumber = sn;
        } else {
            $modal.switchData1[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus4 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.switchData2[index].Status = status;
        if (status === 1) {
            $modal.switchData2[index].SequenceNumber = sn;
        } else {
            $modal.switchData2[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus5 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.switchData3[index].Status = status;
        if (status === 1) {
            $modal.switchData3[index].SequenceNumber = sn;
        } else {
            $modal.switchData3[index].SequenceNumber = -1;
        }
    };
    $scope.setStatus6 = function(index, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        $modal.rangData[index].Status = status;
        if (status === 1) {
            $modal.rangData[index].SequenceNumber = sn;
        } else {
            $modal.rangData[index].SequenceNumber = -1;
        }
    };


    //状态设置
    $scope.setSwitchStatus1 = function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $modal.switchData1[outerIndex].SwitchStatus = status;
    };
    $scope.setSwitchStatus2 = function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $modal.switchData2[outerIndex].SwitchStatus = status;
    };
    $scope.setSwitchStatus3 = function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $modal.switchData3[outerIndex].SwitchStatus = status;
    };

    $modal.openRange = function (currentChannel, number, size, parentSelector) {
        console.log(currentChannel,number);
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-body ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: template.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_channel_config_range.html',
            controller: 'ModalTemplateRangeCtrl',
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
                },
                editing: function () {
                    return $modal.editing;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            console.log("Modal dismissed");
        });
    };



    $modal.ok = function () {
        $modal.templatesChannel = [
            $modal.analogData1,
            $modal.analogData2,
            $modal.switchData1,
            $modal.switchData2,
            $modal.switchData3,
            $modal.rangData
        ];
        console.log($modal.templatesChannel);
        var configUpload = [];

        for(var i = 0; i < $modal.templatesChannel.length; i++){
            for(var j = 0; j < $modal.templatesChannel[i].length; j++){
                var paramId = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Parameter ? $modal.templatesChannel[i][j].Parameter.Id : 0;
                var status = $modal.templatesChannel[i][j] ? $modal.templatesChannel[i][j].Status : 0 ;
                var seqNo = $modal.templatesChannel[i][j] && status === 1 ? $modal.templatesChannel[i][j].SequenceNumber : -1;
                var dataSwitch = $modal.templatesChannel[i][j] ? $modal.templatesChannel[i][j].SwitchStatus : 0 ;
                var ranges = [];
                var chan = i + 1;
                var num = j + 1;
                //功能码
                var fcodeName = 0;
                //MODBUS
                var modbus = 0;
                //位地址
                var bitAddress = 0;
                //高低字节
                var termByte = 0;

                if(paramId===0){
                    continue;
                }

                if(i===0||i===1){
                    chan = i + 1;
                    num = j + 1;
                    fcodeName = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Func ? $modal.templatesChannel[i][j].Func.Id:0;
                    modbus = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Modbus ? $modal.templatesChannel[i][j].Modbus:0;
                    termByte = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Byte? $modal.templatesChannel[i][j].Byte.Id:0 ;
                }
                if(i>=2 && i <5){
                    chan = 3;
                    num = j + (i - 2) * 16 + 1;
                    fcodeName = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Func ? $modal.templatesChannel[i][j].Func.Id:0;
                    modbus = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Modbus? $modal.templatesChannel[i][j].Modbus:0;
                    bitAddress = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].BitAddress? $modal.templatesChannel[i][j].BitAddress:0;
                }
                if(i===5){
                    chan = 5;
                    ranges = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Ranges ? $modal.templatesChannel[i][j].Ranges : [];
                    fcodeName = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Func ? $modal.templatesChannel[i][j].Func.Id:0;
                    modbus = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Modbus ? $modal.templatesChannel[i][j].Modbus:0;
                    termByte = $modal.templatesChannel[i][j] && $modal.templatesChannel[i][j].Byte? $modal.templatesChannel[i][j].Byte.Id:0 ;;
                }




                var configData = {
                    parameter_id: paramId,
                    channel_type: chan,
                    channel_number: num,
                    status: status,
                    sequence_number: seqNo,
                    switch_status: dataSwitch,
                    fcodeName:fcodeName,
                    modbus:parseInt(modbus),
                    termByte:parseInt(termByte),
                    bitAddress:parseInt(bitAddress)
                };

                //表单验证
                if(configData.parameter_id){
                    if(i===0 || i===1 || i===5){
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

                    if(i>=2 && i<5){
                        if(fcodeName===0||modbus===0||bitAddress===0){
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


                if (i === 5 && paramId > 0) {
                    configData.ranges = [];
                    if (ranges.length <= 0) {
                        console.warn("data:", paramId, status, ranges);
                        swal({
                            title: "状态量通道配置错误",
                            text: "已配置的状态量通道，需要完成其状态值的配置才可提交",
                            type: "error"
                        });
                        return;
                    }
                    for (var n in ranges) {
                        var r = ranges[n];
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


        var cParam = {
            // terminal_code:$modal.code,
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
        $http.post("/template_update", {TemplateUpdate:{chan:configUpload,param:cParam,name:$modal.currentData.Name}})
            .then(function (res) {
                App.stopPageLoading();
                swal({
                    title: "模板配置更新成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close('success');
                    currentData = null;
                });
            }, function (err) {
                App.stopPageLoading();
                swal({
                    title: "模板配置更新失败",
                    text: err.data,
                    type: "error"
                });
            });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');
        currentData = null;
    };
});



angular.module('BoilerAdmin').controller('ModalTemplateRangeCtrl', function ($uibModalInstance, $rootScope, $http, $filter, $modal, currentChannel,editing) {
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