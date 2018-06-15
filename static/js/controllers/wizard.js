
//锅炉配置
angular.module('BoilerAdmin').controller("wizardBoilerCtrl",function ($scope,$rootScope,$state ,$stateParams,$http) {

    var uid = $stateParams.uid;

    $scope.getBoiler = function () {
        $http.get("boiler_list/?boiler="+uid).then(function (res) {
            $scope.boiler = res.data[0];
            console.log($scope.boiler);
            $scope.editingCode = false;
            $scope.data = {
                uid: $scope.boiler.Uid,
                name: $scope.boiler.Name,
                registerCode: $scope.boiler.RegisterCode,
                deviceCode: $scope.boiler.DeviceCode,
                factoryNumber: $scope.boiler.FactoryNumber,
                modelCode: $scope.boiler.ModelCode,
                certificateNumber: $scope.boiler.CertificateNumber,

                Usage: "工业锅炉",
                mediumId: $scope.boiler.Medium ? $scope.boiler.Medium.Id : -1,
                fuelId: $scope.boiler.Fuel ? $scope.boiler.Fuel.Uid : "",
                formId: $scope.boiler.Form ? $scope.boiler.Form.Id : 0,
                templateId:$scope.boiler.Template ? $scope.boiler.Template.TemplateId : -1,
                evaporatingCapacity: $scope.boiler.EvaporatingCapacity,

                RegisterOrg: $scope.boiler.RegisterOrg ? $scope.boiler.RegisterOrg : null,
                enterpriseId: $scope.boiler.Enterprise ? $scope.boiler.Enterprise.Uid : "",
                factoryId: $scope.boiler.Factory ? $scope.boiler.Factory.Uid : "",
                maintainerId: $scope.boiler.Maintainer ? $scope.boiler.Maintainer.Uid : "",
                supervisorId: $scope.boiler.Supervisor ? $scope.boiler.Supervisor.Uid : ""
            }
        },function (err) {

        });
    };



    $scope.editingCode = true;

    $scope.links = [];

    $scope.initData = function (uid) {
        if(uid){
            $scope.getBoiler();
        }else {
            var enterprise ;
            var factory;
            if($rootScope.currentUser.Organization && $rootScope.currentUser.Organization.Type.TypeId===1){
                factory = $rootScope.currentUser.Organization.Uid;
            }
            if($rootScope.currentUser.Organization && $rootScope.currentUser.Organization.Type.TypeId===2){
                enterprise = $rootScope.currentUser.Organization.Uid;
            }
            $scope.data = {
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
                templateId: -1,

                evaporatingCapacity: NaN,

                RegisterOrg: null,
                enterpriseId:enterprise? enterprise:"",
                factoryId:factory? factory:"",
                installedId: "",

                links: []
            }
        }

    };

    $scope.init = function () {
        $scope.mediums = [{ Id: -1, Name: '请选择...' }];
        $scope.forms = [{ Id: -1, Name: '请选择...' }];
        $scope.fuels = [{ Uid: '', Name: '请选择...' }];
        $scope.templates = [{ TemplateId: -1, Name: '请选择...' }];

        $scope.enterprises = [{ Uid: '', name: '请选择...' }];
        $scope.factories = [{ Uid: '', name: '请选择...' }];
        $scope.maintainers = [{ Uid: '', name: '请选择...' }];
        $scope.supervisors = [{ Uid: '', name: '请选择...' }];

        $scope.orgTypes = [];

        for (var i in $rootScope.organizations) {
            var org = $rootScope.organizations[i];
            switch (org.Type.TypeId) {
                case 2:
                    if ($scope.enterprises.indexOf(org) === -1) {
                        $scope.enterprises.push(org);
                    }
                    break;
                case 1:
                    if ($scope.factories.indexOf(org) === -1) {
                        $scope.factories.push(org);
                    }
                    break;
                case 5:
                    if ($scope.maintainers.indexOf(org) === -1) {
                        $scope.maintainers.push(org);
                    }
                    break;
                case 7:
                    if ($scope.supervisors.indexOf(org) === -1) {
                        $scope.supervisors.push(org);
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
                    $scope.orgTypes.push(t);
                    break;
            }
        }

        for (var i in $rootScope.boilerMediums) {
            var med = $rootScope.boilerMediums[i];
            if (med.Id === 0 || $scope.mediums.indexOf(med) > -1) {
                continue;
            }

            $scope.mediums.push(med);
        }

        for (var i in $rootScope.boilerForms) {
            var form = $rootScope.boilerForms[i];
            if (form.Id === 0 || $scope.forms.indexOf(form) > -1) {
                continue;
            }

            $scope.forms.push(form);
        }

        console.log($rootScope.boilerTemplates);
        for (var i in $rootScope.boilerTemplates) {
            var template = $rootScope.boilerTemplates[i];
            if (template.TemplateId=== 0 || $scope.templates.indexOf(template) > -1) {
                continue;
            }

            $scope.templates.push(template);
        }
        console.log($scope.templates);

        for (var i in $rootScope.fuels) {
            var fuel = $rootScope.fuels[i];
            if (fuel.Type.Id === 0  || $scope.fuels.indexOf(fuel) > -1) {
                continue;
            }
            $scope.fuels.push(fuel);
        }

        $scope.initData(uid);
    };

    $scope.init();

    $scope.addNewLink = function () {
        if ($scope.links.length >= 4) {
            return;
        }

        $scope.links.push({
            num: $scope.links.length,
        });
    };

    $scope.removeLink = function (link) {
        for (var i in $scope.links) {
            if (link === $scope.links[i]) {
                $scope.links.splice(i, 1);
            }
        }
    };

    $scope.linkTypeChanged = function (link) {
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

    $scope.save = function () {
        console.info("ready to update bInfo!");
        // Ladda.create(document.getElementById('boiler_basic_submit')).start();
        $scope.data.links = [];
        for (var i in $scope.links) {
            var li = $scope.links[i];
            $scope.data.links.push({
                num: li.num,
                type: li.type,
                uid: li.uid
            });
        }


        $http.post("/fast_boiler_add", $scope.data)
            .then(function (res) {
                // console.error("Update bInfo Resp:", res);
                $rootScope.getBoilerList();
                swal({
                    title: "设备基本信息更新成功",
                    type: "success"
                }).then(function () {
                    console.log(res.data);
                    $state.go("wizard.term-bind",{uid:res.data});
                });

                // $scope.initData(res.data);
                // bInfo.currentData = res.data;
                // bInfo.reset();
                // $scope.currentData = res.data;
                // Ladda.create(document.getElementById('boiler_basic_submit')).stop();


            }, function (err) {
                swal({
                    title: "设备基本信息更新失败",
                    text: err.data,
                    type: "error"
                });
                // Ladda.create(document.getElementById('boiler_basic_submit')).stop();
            });
    };

    /*$scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        //currentData = null;
    };*/
});


//终端绑定
angular.module('BoilerAdmin').controller("wizardTermBindCtrl",function ($scope,$rootScope,$state,$stateParams,$http,$uibModal) {

    var uid = $stateParams.uid;
    $scope.terminal = {value:"",bind:false}

    /*$http.get('/terminal_list/?scope=boiler_bind&boiler=' + uid)
        .then(function (res) {
            // $scope.parameters = data;
            var terminals = [];

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
                        if (b.Uid === uid) {
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

            $scope.terminals = terminals;
        });
    */

    $scope.getBoiler = function () {
        $http.get("boiler_list/?boiler="+uid).then(function (res) {
            $scope.boiler = res.data[0];
            if($scope.boiler.TerminalsCombined){
                $scope.terminal.value = $scope.boiler.TerminalsCombined[0].TerminalCode;
                $scope.terminal.bind = true;
            }
            console.log($scope.boiler);
        },function (err) {

        });
    };
    if(uid){
        $scope.getBoiler();
    }


    /*$scope.addTermBind = function (){
        $scope.terminal.push({value:"",bind:false});
    };*/

    $scope.ok = function () {

        $http.post("/fast_terminal_combined", {
            boiler_uid: uid,
            code: parseInt($scope.terminal.value)
        }).then(function (res) {
            // $rootScope.getBoilerList();
            $scope.terminal.bind = true;
            swal({
                title: "绑定设备成功",
                type: "success"
            }).then(function () {
                // currentData = null;
            });
        }, function (err) {
            swal({
                title: "绑定设备失败",
                text: err.data,
                type: "error"
            });
        });
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
        //currentData = null;
    };

     $scope.unbind = function () {
         $http.post("/fast_terminal_unbind", {
             boiler_uid: uid,
             code: parseInt($scope.terminal.value)
         }).then(function (res) {
             // $scope.terminal.splice(index,1);
             $scope.terminal.bind = false;
             swal({
                 title: "绑定已解除",
                 text: "",
                 type: "success"
             });
             // $rootScope.getBoilerList();
         }, function (err) {
             swal({
                 title: "解除绑定失败",
                 text: err.data,
                 type: "error"
             });
         });
    };


     $scope.template = function () {

         var modalInstance = $uibModal.open({
             templateUrl: '/views/wizard/modal-temp.html',
             controller: 'ModalTempSetCtrl',
             size: "",
             windowClass: 'zindex_sub',
             resolve: {
                 code: function () {
                     return $scope.terminal.value;
                 }
             }
         });


         modalInstance.result.then(function (info) {
             $state.go("monitor.thumb");
         }, function () {

         });

     };

    $scope.back = function () {
        $state.go("wizard.boiler",{uid:uid});
    };
    $scope.goNext = function () {
        // $rootScope.getBoilerList();
        $state.go("wizard.term-config",{uid:uid});
    }
    
});

//导入模板
angular.module('BoilerAdmin').controller('ModalTempSetCtrl', function ($scope, $uibModalInstance,$http, code) {

    $scope.code = code;
    $scope.template = "";
    $http.get("/template_list").then(function (res) {
        $scope.templates = res.data;
        // console.log($scope.templates);
    });

    $scope.changeTemp = function (temp) {
        $scope.template = temp;
    };

    $scope.ok = function () {

        $http.post("/fast_template_tocur",{code: parseInt($scope.code),template:$scope.template.Uid})
            .then(function (res) {
                swal({
                    title: "模板导入成功",
                    text: "",
                    type: "success"
                });
                $uibModalInstance.close('success');

            },function (err) {
                swal({
                    title: "模板导入失败",
                    text: err.data,
                    type: "error"
                });
            });

    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});




//终端配置
angular.module('BoilerAdmin').controller("wizardTermConfCtrl",function ($scope,$rootScope,$state,$stateParams,$http,$uibModal) {
    var uid = $stateParams.uid;
    $scope.editing = true;
    $scope.editingCode = true;



    $scope.getBoiler = function () {
        $http.get("boiler_list/?boiler="+uid).then(function (res) {
            $scope.boiler = res.data[0];
            console.log($scope.boiler);
            $scope.terminal = $scope.boiler.TerminalsCombined[0];
            console.log($scope.terminal);
        },function (err) {

        });
    };
    if(uid){
        $scope.getBoiler();
    }




    //功能码
    $http.get("/term_function_code_list").then(function (res) {
        $scope.fcode = res.data;
        $scope.fcode2 = [
            {Id: 1, Name: "01", Value: 1},
            {Id: 2, Name: "02", Value: 2},
            {Id: 3, Name: "03", Value: 3},
            {Id: 99, Name: "None", Value: 99}
        ];
        $scope.fcode1 = [
            {Id: 3, Name: "03", Value: 3},
            {Id: 4, Name: "04", Value: 4}

        ];
    });

    //高低字节
    $http.get("/term_byte_list").then(function (res) {
        $scope.hlCodes = res.data;
    });

    //通信接口地址
    $http.get("/correspond_type_list").then(function (res) {
        $scope.communiInterfaces = res.data;
    });
    //数据位
    $http.get("/date_bit_list").then(function (res) {
        $scope.dataBits = res.data;
    });
    //心跳包频率
    $http.get("/heartbeat_packet_list").then(function (res) {
        $scope.heartbeats = res.data;
    });
    //校验位
    $http.get("/parity_bit").then(function (res) {
        $scope.checkDigits = res.data;
    });
    //从机地址
    $http.get("/slave_address_list").then(function (res) {
        $scope.subAdrs = res.data;
    });
    //停止位
    $http.get("/stop_bit_list").then(function (res) {
        $scope.stopBits = res.data;
    });
    //波特率
    $http.get("/baud_rate_list").then(function (res) {
        $scope.BaudRates = res.data;
    });

    $scope.priorities=[];
    for(var i = 0; i<48; i++){
        $scope.priorities.push(i);
    }



    //模拟通道
    var aNum = 1;
    $scope.analogueList = [
        {
            ChannelNumber:aNum,
            Parameter:{
                Name:"",
                Scale:null,
                Unit:""
            },
            Func:null,
            Byte:null,
            Modbus:null,
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        }
    ];

    //开关通道
    var sNum = 2;
    $scope.switchList = [
        {
            ChannelNumber:1,
            Parameter:{
                Name:"点火信号",
            },
            Func:99,
            Modbus:null,
            BitAddress:null,
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        },
        {
            ChannelNumber:2,
            Parameter:{
                Name:"PLC",
            },
            Func:99,
            Modbus:null,
            BitAddress:null,
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        }
    ];

    //状态通道
    var rNum = 1;
    $scope.rangeList = [
        {
            ChannelNumber:rNum,
            Parameter:{
                Name:"",
            },
            Func:null,
            Byte:null,
            Modbus:null,
            Ranges: [],
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        }
    ];

    $scope.infomation = {
        BaudRate:null,
        dataBit:null,
        stopBit:null,
        checkDigit:null,
        communiInterface:null,
        subAdr:null,
        heartbeat:null,
    };




    $scope.dataChanged = function (data) {
        if (!data.Parameter.Name) {
            data.Parameter = {};
            data.Status = -1;
            data.SwitchStatus = 0;
            data.Ranges = null;
        } else {
            if (!data.Status || data.Status === -1) {
                data.Status = 0;
            }

            if (!data.SwitchStatus || data.SwitchStatus === 0) {
                data.SwitchStatus = 1;
            }

            // console.log(data);
        }

    };



    //添加
    $scope.addAnalogue = function () {
        aNum++;
        $scope.analogueList.push({
            ChannelNumber:aNum,
            Parameter:{
                Name:"",
                Scale:null,
                Unit:""
            },
            Func:null,
            Byte:null,
            Modbus:null,
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        });
    };
    $scope.addSwitch = function () {
        sNum++;
        $scope.switchList.push({
            ChannelNumber:sNum,
            Parameter:{
                Name:"",
            },
            Func:null,
            Modbus:null,
            BitAddress:null,
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        });
    };
    $scope.addRange = function () {
        rNum++;
        $scope.rangeList.push({
            ChannelNumber:rNum,
            Parameter:{
                Name:"",
            },
            Func:null,
            Byte:null,
            Modbus:null,
            Ranges: [],
            SequenceNumber: 0,
            Status: 0,
            SwitchStatus: 0
        });
    };

    //删除
    $scope.removeAnalogue = function (index) {
        $scope.analogueList.splice(index,1);
    };
    $scope.removeSwitch = function (index) {
        $scope.switchList.splice(index,1);
    };
    $scope.removeRange = function (index) {
        $scope.rangeList.splice(index,1);
    };



    $scope.fCodeChange =function (data) {
        if(data.Func ===1||data.Func ===2){
            data.BitAddress = 1;
        }
    };


    //位置设置
    $scope.setStatus = function(data, status, sn) {
        // console.warn("$scope.setStatus", index, status, sn);
        data.Status = status;
        if (status === 1) {
            data.SequenceNumber = sn;
        } else {
            data.SequenceNumber = -1;
        }
    };





    //状态设置
    $scope.setSwitchStatus= function(outerIndex, status) {
        console.warn("$scope.setSwitchStatus", outerIndex, status);
        $scope.switchList[outerIndex].SwitchStatus = status;
        console.log($scope.switchList);
    };


    $scope.openRange = function (currentChannel, number, size) {
        console.log(currentChannel,number);
        var modalInstance = $uibModal.open({
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/terminal_channel_config_range.html',
            controller: 'ModalWizardRangeCtrl',
            controllerAs: '$modalRange',
            size: size,
            windowClass: 'zindex_sub',
            resolve: {
                $modal: function () {
                    return $scope;
                },
                currentChannel: function () {
                    currentChannel.ChannelNumber = number;
                    return currentChannel;
                },
                editing: function () {
                    return $scope.editing;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            terminal.selected = selectedItem;
        }, function () {
            console.log("Modal dismissed");
        });
    };


    $scope.back = function () {
        $state.go("wizard.term-bind",{uid:uid});
    };



    //提交
    $scope.ok = function () {

        //模拟通道
        var analogueList = [];
        var aNumList = [];
        for(var i =0; i<$scope.analogueList.length;i++){
            if($scope.analogueList[i].Parameter.Name){
                if($scope.analogueList[i].ChannelNumber>24){
                    swal({
                        title: "通道配置更新失败",
                        text:"模拟通道不能超过24",
                        type: "error"
                    });
                    return false;
                }
                aNumList.push($scope.analogueList[i].ChannelNumber);
                if(!$scope.analogueList[i].Func||!$scope.analogueList[i].Byte||!$scope.analogueList[i].Modbus|| !$scope.analogueList[i].Parameter.Scale){
                    swal({
                        title: "通道配置更新失败",
                        text:"模拟通道配置信息不全 ，参数不能为0 "+ i,
                        type: "error"
                    });
                    return false;
                }
                if($scope.analogueList[i].Func===3){
                    if($scope.analogueList[i].Modbus<=40000||$scope.analogueList[i].Modbus>=50000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为03，MODBUS地址范围40001-49999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.analogueList[i].Func===4){
                    if($scope.analogueList[i].Modbus<=30000||$scope.analogueList[i].Modbus>=40000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为04，MODBUS地址范围30001-39999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                analogueList.push($scope.analogueList[i]);
            }
        }
        aNumList.sort();
        console.log(aNumList);
        for(var i = 0; i<aNumList.length; i++){
            if(aNumList[i]===aNumList[i+1]){
                swal({
                    title: "通道重复配置",
                    text:" 模拟量通道数不能相同",
                    type: "error"
                });
                return false;
            }
        }

        //开关通道
        var switchList = [];
        var sNumList = [];
        for(var i =0; i<$scope.switchList.length;i++){
            if($scope.switchList[i].Parameter.Name){
                if($scope.switchList[i].ChannelNumber>48){
                    swal({
                        title: "通道配置更新失败",
                        text:"开关通道不能超过48",
                        type: "error"
                    });
                    return false;
                }
                aNumList.push($scope.switchList[i].ChannelNumber);
                if($scope.switchList[i].Func!==99 &&(!$scope.switchList[i].Func||!$scope.switchList[i].BitAddress||!$scope.switchList[i].Modbus)){
                    swal({
                        title: "通道配置更新失败",
                        text:"开关通道配置信息不全 ，参数不能为0 "+ i,
                        type: "error"
                    });
                    return false;
                }
                if($scope.switchList[i].Func===1){
                    if($scope.switchList[i].Modbus<1||$scope.switchList[i].Modbus>=10000){
                        swal({
                            title: "开关通道MODBUS地址错误",
                            text:"功能码为01，MODBUS地址范围00001-09999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                    if($scope.switchList[i].BitAddress!=1){
                        swal({
                            title: "开关通道位地址错误",
                            text:"功能码为01，对应位地址为1" ,
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.switchList[i].Func===2){
                    if($scope.switchList[i].Modbus<=10000||$scope.switchList[i].Modbus>=20000){
                        swal({
                            title: "开关通道MODBUS地址错误",
                            text:"功能码为02，MODBUS地址范围10001-19999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                    if($scope.switchList[i].BitAddress!=1){
                        swal({
                            title: "开关通道位地址错误",
                            text:"功能码为02，对应位地址为1",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.switchList[i].Func===3){
                    if($scope.switchList[i].Modbus<=40000||$scope.switchList[i].Modbus>=50000){
                        swal({
                            title: "开关通道MODBUS地址错，请修改",
                            text:"功能码为03，MODBUS地址范围40001-49999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                    if($scope.switchList[i].BitAddress<1||$scope.switchList[i].BitAddress>16){
                        swal({
                            title: "位地址错误",
                            text: "功能码为03，对应位地址范围为1-16",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }

                switchList.push($scope.switchList[i]);
            }
        }
        sNumList.sort();
        for(var i = 0; i<sNumList.length; i++){
            if(sNumList[i]===sNumList[i+1]){
                swal({
                    title: "通道重复配置",
                    text:" 开关量通道数不能相同",
                    type: "error"
                });
                return false;
            }
        }

        //状态通道
        var rangeList = [];
        var rNumList = [];
        for(var i =0; i<$scope.rangeList.length;i++){
            var range = angular.copy($scope.rangeList[i].Ranges);
            console.log(range);
            if($scope.rangeList[i].Parameter.Name){
                if($scope.rangeList[i].ChannelNumber>12){
                    swal({
                        title: "通道配置更新失败",
                        text:"状态通道不能超过12",
                        type: "error"
                    });
                    return false;
                }
                aNumList.push($scope.rangeList[i].ChannelNumber);
                if(!$scope.rangeList[i].Func||!$scope.rangeList[i].Byte||!$scope.rangeList[i].Modbus){
                    swal({
                        title: "通道配置更新失败",
                        text:"模拟通道配置信息不全 ，参数不能为0 "+ i,
                        type: "error"
                    });
                    return false;
                }
                if($scope.rangeList[i].Func===3){
                    if($scope.rangeList[i].Modbus<=40000||$scope.rangeList[i].Modbus>=50000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为03，MODBUS地址范围40001-49999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.rangeList[i].Func===4){
                    if($scope.rangeList[i].Modbus<=30000||$scope.rangeList[i].Modbus>=40000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为04，MODBUS地址范围30001-39999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }

                $scope.rangeList[i].Ranges = [];
                if (range.length <= 0) {
                    swal({
                        title: "状态量通道配置错误",
                        text: "已配置的状态量通道，需要完成其状态值的配置才可提交",
                        type: "error"
                    });
                    return;
                }
                for (var n in range) {
                    var r = range[n];
                    var rg = {};
                    rg.Min = r.Min;
                    rg.Max = r.Max;
                    rg.Name = r.Name;
                    switch (typeof n) {
                        case "number":
                            rg.Value = n;
                            break;
                        case "string":
                            rg.Value = parseInt(n, 10);
                            break;
                    }

                    $scope.rangeList[i].Ranges.push(rg);
                }

                rangeList.push($scope.rangeList[i]);
            }
        }
        rNumList.sort();
        console.log(rNumList);
        for(var i = 0; i<rNumList.length; i++){
            if(rNumList[i]===rNumList[i+1]){
                swal({
                    title: "通道重复配置",
                    text:" 状态量通道数不能相同",
                    type: "error"
                });
                return false;
            }
        }




        $scope.channel ={
            Analogue:analogueList,
            Switch:switchList,
            Range:rangeList
        };


        //通信参数
        var cParam = {
            baudRate : $scope.infomation.BaudRate?$scope.infomation.BaudRate.Id:0,
            dataBit : $scope.infomation.dataBit?$scope.infomation.dataBit.Id:0,
            stopBit : $scope.infomation.stopBit?$scope.infomation.stopBit.Id:0,
            checkDigit : $scope.infomation.checkDigit?$scope.infomation.checkDigit.Id:0,
            communInterface : $scope.infomation.communiInterface?$scope.infomation.communiInterface.Id:0,
            slaveAddress : $scope.infomation.subAdr?$scope.infomation.subAdr.Id:0,
            heartbeat:$scope.infomation.heartbeat?$scope.infomation.heartbeat.Id:0
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

        console.log($scope.channel,cParam);
        $scope.cParam = cParam;



        if($scope.terminalPass!=="1234567"){
            swal({
                title: "终端密码错误",
                text:" ",
                type: "error"
            });
            return false;
        }

        $http.post("/fast_term_channel_config",{Chan:$scope.channel, Param:$scope.cParam, Code:$scope.terminal.TerminalCode})
            .then(function (res) {
                swal({
                    title: "终端配置更新成功",
                    type: "success"
                });
                $state.go("monitor.thumb");
            },function (err) {
                swal({
                    title: "终端配置更新失败",
                    text: err.data,
                    type: "error"
                });
            })


    }




//    模板另存为
    $scope.templateSave = function () {
        //模拟通道
        var analogueList = [];
        var aNumList = [];
        for(var i =0; i<$scope.analogueList.length;i++){
            if($scope.analogueList[i].Parameter.Name){
                if($scope.analogueList[i].ChannelNumber>24){
                    swal({
                        title: "通道配置更新失败",
                        text:"模拟通道不能超过24",
                        type: "error"
                    });
                    return false;
                }
                aNumList.push($scope.analogueList[i].ChannelNumber);
                if(!$scope.analogueList[i].Func||!$scope.analogueList[i].Byte||!$scope.analogueList[i].Modbus|| !$scope.analogueList[i].Parameter.Scale){
                    swal({
                        title: "通道配置更新失败",
                        text:"模拟通道配置信息不全 ，参数不能为0 "+ i,
                        type: "error"
                    });
                    return false;
                }
                if($scope.analogueList[i].Func===3){
                    if($scope.analogueList[i].Modbus<=40000||$scope.analogueList[i].Modbus>=50000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为03，MODBUS地址范围40001-49999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.analogueList[i].Func===4){
                    if($scope.analogueList[i].Modbus<=30000||$scope.analogueList[i].Modbus>=40000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为04，MODBUS地址范围30001-39999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                analogueList.push($scope.analogueList[i]);
            }
        }
        aNumList.sort();
        console.log(aNumList);
        for(var i = 0; i<aNumList.length; i++){
            if(aNumList[i]===aNumList[i+1]){
                swal({
                    title: "通道重复配置",
                    text:" 模拟量通道数不能相同",
                    type: "error"
                });
                return false;
            }
        }

        //开关通道
        var switchList = [];
        var sNumList = [];
        for(var i =0; i<$scope.switchList.length;i++){
            if($scope.switchList[i].Parameter.Name){
                if($scope.switchList[i].ChannelNumber>48){
                    swal({
                        title: "通道配置更新失败",
                        text:"开关通道不能超过48",
                        type: "error"
                    });
                    return false;
                }
                aNumList.push($scope.switchList[i].ChannelNumber);
                if($scope.switchList[i].Func!==99 &&(!$scope.switchList[i].Func||!$scope.switchList[i].BitAddress||!$scope.switchList[i].Modbus)){
                    swal({
                        title: "通道配置更新失败",
                        text:"开关通道配置信息不全 ，参数不能为0 "+ i,
                        type: "error"
                    });
                    return false;
                }
                if($scope.switchList[i].Func===1){
                    if($scope.switchList[i].Modbus<1||$scope.switchList[i].Modbus>=10000){
                        swal({
                            title: "开关通道MODBUS地址错误",
                            text:"功能码为01，MODBUS地址范围00001-09999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                    if($scope.switchList[i].BitAddress!=1){
                        swal({
                            title: "开关通道位地址错误",
                            text:"功能码为01，对应位地址为1" ,
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.switchList[i].Func===2){
                    if($scope.switchList[i].Modbus<=10000||$scope.switchList[i].Modbus>=20000){
                        swal({
                            title: "开关通道MODBUS地址错误",
                            text:"功能码为02，MODBUS地址范围10001-19999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                    if($scope.switchList[i].BitAddress!=1){
                        swal({
                            title: "开关通道位地址错误",
                            text:"功能码为02，对应位地址为1",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.switchList[i].Func===3){
                    if($scope.switchList[i].Modbus<=40000||$scope.switchList[i].Modbus>=50000){
                        swal({
                            title: "开关通道MODBUS地址错，请修改",
                            text:"功能码为03，MODBUS地址范围40001-49999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                    if($scope.switchList[i].BitAddress<1||$scope.switchList[i].BitAddress>16){
                        swal({
                            title: "位地址错误",
                            text: "功能码为03，对应位地址范围为1-16",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }

                switchList.push($scope.switchList[i]);
            }
        }
        sNumList.sort();
        for(var i = 0; i<sNumList.length; i++){
            if(sNumList[i]===sNumList[i+1]){
                swal({
                    title: "通道重复配置",
                    text:" 开关量通道数不能相同",
                    type: "error"
                });
                return false;
            }
        }

        //状态通道
        var rangeList = [];
        var rNumList = [];
        for(var i =0; i<$scope.rangeList.length;i++){
            var range = angular.copy($scope.rangeList[i].Ranges);
            console.log(range);
            if($scope.rangeList[i].Parameter.Name){
                if($scope.rangeList[i].ChannelNumber>12){
                    swal({
                        title: "通道配置更新失败",
                        text:"状态通道不能超过12",
                        type: "error"
                    });
                    return false;
                }
                aNumList.push($scope.rangeList[i].ChannelNumber);
                if(!$scope.rangeList[i].Func||!$scope.rangeList[i].Byte||!$scope.rangeList[i].Modbus){
                    swal({
                        title: "通道配置更新失败",
                        text:"模拟通道配置信息不全 ，参数不能为0 "+ i,
                        type: "error"
                    });
                    return false;
                }
                if($scope.rangeList[i].Func===3){
                    if($scope.rangeList[i].Modbus<=40000||$scope.rangeList[i].Modbus>=50000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为03，MODBUS地址范围40001-49999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }
                if($scope.rangeList[i].Func===4){
                    if($scope.rangeList[i].Modbus<=30000||$scope.rangeList[i].Modbus>=40000){
                        swal({
                            title: "MODBUS地址错误",
                            text:"功能码为04，MODBUS地址范围30001-39999",
                            type: "error"
                        });
                        App.stopPageLoading();
                        return false;
                    }
                }

                $scope.rangeList[i].Ranges = [];
                if (range.length <= 0) {
                    swal({
                        title: "状态量通道配置错误",
                        text: "已配置的状态量通道，需要完成其状态值的配置才可提交",
                        type: "error"
                    });
                    return;
                }
                for (var n in range) {
                    var r = range[n];
                    var rg = {};
                    rg.Min = r.Min;
                    rg.Max = r.Max;
                    rg.Name = r.Name;
                    switch (typeof n) {
                        case "number":
                            rg.Value = n;
                            break;
                        case "string":
                            rg.Value = parseInt(n, 10);
                            break;
                    }

                    $scope.rangeList[i].Ranges.push(rg);
                }

                rangeList.push($scope.rangeList[i]);
            }
        }
        rNumList.sort();
        console.log(rNumList);
        for(var i = 0; i<rNumList.length; i++){
            if(rNumList[i]===rNumList[i+1]){
                swal({
                    title: "通道重复配置",
                    text:" 状态量通道数不能相同",
                    type: "error"
                });
                return false;
            }
        }




        $scope.channel ={
            Analogue:analogueList,
            Switch:switchList,
            Range:rangeList
        };


        //通信参数
        var cParam = {
            baudRate : $scope.infomation.BaudRate?$scope.infomation.BaudRate.Id:0,
            dataBit : $scope.infomation.dataBit?$scope.infomation.dataBit.Id:0,
            stopBit : $scope.infomation.stopBit?$scope.infomation.stopBit.Id:0,
            checkDigit : $scope.infomation.checkDigit?$scope.infomation.checkDigit.Id:0,
            communInterface : $scope.infomation.communiInterface?$scope.infomation.communiInterface.Id:0,
            slaveAddress : $scope.infomation.subAdr?$scope.infomation.subAdr.Id:0,
            heartbeat:$scope.infomation.heartbeat?$scope.infomation.heartbeat.Id:0
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

        console.log($scope.channel,cParam);
        $scope.cParam = cParam;



        if($scope.terminalPass!=="1234567"){
            swal({
                title: "终端密码错误",
                text:" ",
                type: "error"
            });
            return false;
        }

        var modalInstance = $uibModal.open({
            templateUrl: '/directives/modal/terminal_channel_template.html',
            controller: 'ModalTemplateSaveCtrl',
            size: " ",
            windowClass: 'zindex_sub',
            resolve: {
                cParam: function () {
                    return $scope.cParam;
                },
                configUpload: function () {
                    return $scope.channel;
                }
            }
        });


        modalInstance.result.then(function (selectedItem) {
            // $scope.selected = selectedItem;
        }, function () {

        });
    }




});


//模板另存为
angular.module('BoilerAdmin').controller('ModalTemplateSaveCtrl', function ($scope, $http, $uibModalInstance,cParam,configUpload) {

    console.log("cParam:",cParam,"configUpload:",configUpload);

    $scope.templateName = "";
    $scope.ok = function () {

        if(!$scope.templateName){
            swal({
                title: "模板未命名",
                text: "请填写需要保存的模板名称" ,
                type: "error"
            });
            return;
        }
        App.startPageLoading({message: '正在加载数据...'});
        $http.post("/fast_template", {Name:$scope.templateName,Chan:configUpload,Param:cParam})
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

    };


    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});




angular.module('BoilerAdmin').controller('ModalWizardRangeCtrl', function ($uibModalInstance, $rootScope, $http, $filter, $scope, currentChannel,editing) {
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
        console.log($modalRange.ranges);
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




