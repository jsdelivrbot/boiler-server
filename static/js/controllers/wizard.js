
//锅炉配置
angular.module('BoilerAdmin').controller("wizardBoilerCtrl",function ($scope,$rootScope,$state ,$stateParams,$http) {

    var uid = $stateParams.uid;
    if(!uid){
        $scope.currentData = null;
    }
    $scope.editingCode = true;

    $scope.links = [];

    $scope.initData = function (uid) {
        if(uid){

        }else {
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
                // enterpriseId: "",
                // factoryId: "",
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
angular.module('BoilerAdmin').controller("wizardTermBindCtrl",function ($scope,$rootScope,$state,$stateParams,$http) {

    var uid = $stateParams.uid;
    /*$scope.currentData = currentData;
    $scope.name = currentData.Name;*/
    $scope.terminal = [
        {value:"",bind:false}
    ];

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
            console.log($scope.boiler);
        },function (err) {

        });
    };
    $scope.getBoiler();


    $scope.addTermBind = function (){
        $scope.terminal.push({value:"",bind:false});
    };

    $scope.ok = function (term) {
        // console.info("ready to bind boiler!");
        $http.post("/fast_terminal_combined", {
            boiler_uid: uid,
            code: parseInt(term.value)
        }).then(function (res) {
            // $rootScope.getBoilerList();
            term.bind = true;
            // $scope.getBoiler();
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

     $scope.unbind = function (index) {
         $scope.terminal.splice(index,1);
         $http.post("/boiler_unbind/", {
             boiler_id: uid,
             terminal_id: term[index].value
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
    };
    
    $scope.goNext = function () {
        // $rootScope.getBoilerList();
        $state.go("wizard.term-config",{uid:uid});
    }
    
});



//终端配置
angular.module('BoilerAdmin').controller("wizardTermConfCtrl",function ($scope,$rootScope,$state,$stateParams,$http) {

});










