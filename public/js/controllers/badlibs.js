//test
define(["app", "services/badlibs"],function(app){


	app.controller("badlibsController", function($scope, badlibsFactory, $rootScope){
		$scope.test = "Hello";
		$scope.lib = {};
		$scope.partOfSpeech = {};
		$scope.rating = {};

		$scope.allLibs = badlibsFactory.find($scope.lib).
			then(function(data){
				$scope.allLibs = data;
			},function(err){
				$scope.err = err;
			});

		$scope.allPartsOfSpeech =  badlibsFactory.findPartsOfSpeech($scope.partOfSpeech).
			then(function(data){
				$scope.allPartsOfSpeech = data;
			},function(err){
				$scope.err = err;
			});

		$scope.allRatings =  badlibsFactory.findRatings($scope.rating).
			then(function(data){
				$scope.allRatings = data;
			},function(err){
				$scope.err = err;
			});

		$scope.createLib = function(){
			//TODO
		}

		$scope.append = function(partOfSpeech){
			console.log(partOfSpeech);
			//TODO
			//http://stackoverflow.com/questions/1064089/inserting-a-text-where-cursor-is-using-javascript-jquery
		}

	});
});