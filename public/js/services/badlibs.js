define(['app'], function (app) {
	app.factory("badlibsFactory",function($q, $http, $rootScope){
		var factory = {};

		factory.find = function(lib){
			var deferred = $q.defer();
			$http({
				method:'post',
				url:'/lib/find',
				data: lib
			}).success(function(data){
				deferred.resolve(data)
			}).error(function(data){
				deferred.reject(data);
			});
			return deferred.promise;
		};

		factory.findPartsOfSpeech = function(partOfSpeech){
			var deferred = $q.defer();
			$http({
				method:'post',
				url:'/partsofspeech/find',
				data: partOfSpeech
			}).success(function(data){
				deferred.resolve(data)
			}).error(function(data){
				deferred.reject(data);
			});
			return deferred.promise;
		};

		factory.findRatings = function(rating){
			var deferred = $q.defer();
			$http({
				method:'post',
				url:'/ratings/find',
				data: rating
			}).success(function(data){
				deferred.resolve(data)
			}).error(function(data){
				deferred.reject(data);
			});
			return deferred.promise;
		};

		return factory;
	});
});