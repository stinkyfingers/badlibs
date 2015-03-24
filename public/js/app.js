//app
define(['angularAMD', 'angular-route', 'angularCookies'], function (angularAMD) {


	var  app = angular.module("app",["ngRoute","ngCookies"],function($interpolateProvider){
		$interpolateProvider.startSymbol('[[');
	    $interpolateProvider.endSymbol(']]');
	});

	app.config(function($routeProvider, $locationProvider){

		if(window.history && window.history.pushState){
		      $locationProvider.html5Mode({
				  enabled: true,
				  requireBase: false
				});
		    }

		$routeProvider.
			
			when("/create",angularAMD.route({
				templateUrl: '/public/js/templates/createbadlibs.tmpl',
				controller: 'badlibsController',
				controllerUrl: 'controllers/badlibs'
			})).
			when("/badlibs/:id",angularAMD.route({
				templateUrl: '/public/js/templates/playbadlibs.tmpl',
				controller: 'badlibsController',
				controllerUrl: 'controllers/badlibs'
			})).
			when("/",angularAMD.route({
				templateUrl: '/public/js/templates/badlibs.tmpl',
				controller: 'badlibsController',
				controllerUrl: 'controllers/badlibs'
			})).
			
			otherwise({redirectTo: "/"});

		

	});

	app.run(function($rootScope, $cookies){
		var u = $cookies.user;
		$rootScope.user = u;
	});



  return angularAMD.bootstrap(app);
});