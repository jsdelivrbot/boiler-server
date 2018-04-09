angular.module('BoilerAdmin').controller("templateCtrl",function ($rootScope,$scope,$uibModal,$http,settings,DTOptionsBuilder, DTColumnDefBuilder) {
    var template = this;

    App.initAjax();

    // dialogue.refreshDataTables();

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

    template.datasource=[
        {
            num:1,
            name:"通用模板一",
            enterprise:"锅炉制造厂"
        },
        {
            num:2,
            name:"通用模板二",
            enterprise:"二号锅炉制造厂"
        }
    ];


    var currentData;
    var editing;

    template.new = function () {
        currentData = null;
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
            title: "确认删除该文件？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            // $http.post("/dialogue_delete/", {
            //     uid: uid
            // }).then(function (res) {
            //     swal({
            //         title: "文件删除成功",
            //         type: "success"
            //     }).then(function () {
            //         dialogue.refreshDataTables();
            //     });
            // }, function (err) {
            //     swal({
            //         title: "删除文件失败",
            //         text: err.data,
            //         type: "error"
            //     });
            // });
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



angular.module('BoilerAdmin').controller('ModalEditTemplateCtrl', function ($rootScope,$scope, $uibModalInstance,$http,currentData,editing) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.editing = editing;
    $modal.editingCode = true;
    // $modal.template = true;

    $modal.category = 9;

    //下发test
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






    $modal.priorities = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];

    if(currentData==null){
        currentData = {
            code:null
        };
    }
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

                if($modal.chanMatrix[i][j].Analogue.Modbus===0){
                    $modal.chanMatrix[i][j].Analogue.Modbus = null;
                }
                if($modal.chanMatrix[i][j].Switch.Modbus===0){
                    $modal.chanMatrix[i][j].Switch.Modbus = null;
                }
                if($modal.chanMatrix[i][j].Switch.BitAddress===0){
                    $modal.chanMatrix[i][j].Switch.BitAddress = null;
                }



                if (!$modal.chanMatrix[i][j].RuntimeParameterChannelConfig) {
                    $modal.chanMatrix[i][j] = {
                        Name: "默认(未配置)",
                        noStatus:true
                    }

                }

                if ((i !== 0 ||j !== 2 ) &&  (!$modal.dataMatrix[i][j].RuntimeParameterChannelConfig || $modal.dataMatrix[i][j].RuntimeParameterChannelConfig.IsDefault) ) {
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

    $modal.matrixChanged = function (outerIndex, innerIndex) {
        console.info("Data Matrix:", $modal.dataMatrix, "\n", outerIndex, ":", innerIndex);
        if ($modal.dataMatrix[outerIndex][innerIndex].Parameter.Id === 0) {
            $modal.dataMatrix[outerIndex][innerIndex].Parameter = null;
            $modal.dataMatrix[outerIndex][innerIndex].oParamId = 0;
            $modal.dataMatrix[outerIndex][innerIndex].Status = -1;
            $modal.dataMatrix[outerIndex][innerIndex].SwitchStatus = 0;
            $modal.dataMatrix[outerIndex][innerIndex].Ranges = null;
            if($modal.chanMatrix[outerIndex][innerIndex].IsDefault!==true){
                $modal.chanMatrix[outerIndex][innerIndex].Name="默认(未配置)"
            }
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
            console.log($modal.chanMatrix);
        }

    };

    //恢复默认
    $scope.matrixReset = function () {
        for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                $modal.dataMatrix[i][j] = null;
                if($modal.chanMatrix[i][j].IsDefault!=true){
                    $modal.chanMatrix[i][j].Name="默认(未配置)"
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
        if(fcode.Id ===1||fcode.Id ===2){
            $modal.chanMatrix[i][j].Switch.BitAddress = 1;
        }
    };


    //位置设置
    $scope.setStatus = function(outerIndex, innerIndex, status, sn) {
        // console.warn("$scope.setStatus", outerIndex, innerIndex, status, sn);
        $modal.dataMatrix[outerIndex][innerIndex].Status = status;
        if (status === 1) {
            $modal.dataMatrix[outerIndex][innerIndex].SequenceNumber = sn;
        } else {
            $modal.dataMatrix[outerIndex][innerIndex].SequenceNumber = -1;
        }
    };

    //状态设置
    $scope.setSwitchStatus = function(outerIndex, innerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, innerIndex, status);
        $modal.dataMatrix[outerIndex][innerIndex].SwitchStatus = status;
    };


    $modal.ok = function () {
        if (!$modal.code.length || $modal.code.length !== 6) {
            console.error("$modal.code error:", $modal.code);
            return;
        }
        // Ladda.create(document.getElementById('channel_ok')).start();

        var configUpload = [];
        for (var i = 0; i < $modal.dataMatrix.length; i++) {
            for (var j = 0; j < $modal.dataMatrix[i].length; j++) {
                if ($modal.dataMatrix[i][j] !== $modal.chanMatrix[i][j]) {
                    if ((!$modal.dataMatrix[i][j] /*|| !$modal.dataMatrix[i][j].Parameter*/) && ($modal.chanMatrix[i][j] && $modal.chanMatrix[i][j].IsDefault === true)) {
                        console.warn('!!NULL data:', $modal.dataMatrix[i][j], $modal.chanMatrix[i][j]);
                        continue;
                    }
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

                            switch_status: dataSwitch
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


        console.warn("$modal channel update!", configUpload);

        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/channel_config_update/", configUpload)
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
            });
        // Ladda.create(document.getElementById('channel_ok')).stop();
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');
        currentData = null;
    };
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
