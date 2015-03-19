//test
define(["app", "services/badlibs"],function(app){


	app.controller("badlibsController", function($scope, badlibsFactory, $rootScope){
		$scope.test = "Hello";
		$scope.lib = {};

		$scope.allLibs = badlibsFactory.find($scope.lib).
			then(function(data){
				$scope.allLibs = data;
			},function(err){
				$scope.err = err;
			});

		$scope.createLib = function{
			
		}

	});
});