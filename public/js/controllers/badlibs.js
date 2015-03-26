//test
define(["app", "services/badlibs"],function(app){


	app.controller("badlibsController", function($scope, badlibsFactory, $rootScope, $route){
		$scope.test = "Hello";
		$scope.lib = {};
		$scope.partOfSpeech = {};
		$scope.rating = {};
		$scope.posArray = [];
		$scope.pos = {};
		$scope.objArray = [];
		// $scope.completedLib = {};
		var iterator = 0;


		params = $route.current.params;

		if ((typeof(params.id) != "undefined") && (params.id != "")){
			$scope.lib._id = params.id;
			badlibsFactory.find($scope.lib).then(function(data){
				$scope.lib = data[0];
				setup($scope.lib);
			},function(err){
				$scope.err = err;
			});
		}

		setup = function(lib){
			var re = /\(\([A-Za-z \-]*\)\)/gi;
			$scope.posArray = lib.text.match(re);
			$scope.objArray = new Array($scope.posArray.length);
			angular.forEach($scope.posArray,function(v,k){
				var pos = v.replace("((","");
				pos = pos.replace("))","");
				$scope.posArray[k] = pos;
				$scope.objArray[k] = {};
				$scope.objArray[k].pos = pos;

			})
		}

		$scope.submitPos = function(){
			//TODO - required
			// var ok =$error.required;
			// console.log(ok)
			var re =  /\(\([A-Za-z \-]*\)\)/gi;
			var newString = $scope.lib.text.replace(re, replacer);
			iterator = 0;
			//check capitalization
			// newString = newString.replace(/[!\?\.][ ][a-z]/,function(){return arguments[0].toUpperCase();})

			$scope.completedLib = newString;
		}

		replacer = function(){
			var output;
			angular.forEach($scope.objArray,function(v,k){
				if (k == iterator){
					output =  v.userWord;
				}
			});

			iterator++;
			return output;
		}

		$scope.allLibs = badlibsFactory.find($scope.lib).
			then(function(data){
				$scope.allLibs = data;
				$scope.libsView = $scope.allLibs;
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
		$scope.clear = function(){
			$scope.lib = {};
		}


		$scope.createLib = function(lib){
			if (lib.rating =="" || typeof(lib.rating) == "undefined"){
				lib.rating = "R";
			}
			badlibsFactory.createLib(lib).then(function(data){
				$scope.lib = data;
				alert("Success!");
				$scope.lib = {};
			},function(err){
				$scope.err = err;
			});
		}

		$scope.append = function(partOfSpeech){
			textarea = document.querySelector("#libtext");
			//position of cursor
			caret = getCaret(textarea);

			beginning = textarea.value.substr(0,caret);
			end = textarea.value.substr(caret, textarea.value.length);

			$scope.lib.text = textarea.value = beginning + ("(("+ partOfSpeech.value+"))") + end;

			textarea.focus();
		}


		//find position of carets
		function getCaret(el) {
		  if (el.selectionStart) { 
		    return el.selectionStart; 
		  } else if (document.selection) { 
		    el.focus(); 

		    var r = document.selection.createRange(); 
		    if (r == null) { 
		      return 0; 
		    } 

		    var re = el.createTextRange(), 
		    rc = re.duplicate(); 
		    re.moveToBookmark(r.getBookmark()); 
		    rc.setEndPoint('EndToStart', re); 

		    var add_newlines = 0;
		    for (var i=0; i<rc.text.length; i++) {
		      if (rc.text.substr(i, 2) == '\r\n') {
		        add_newlines += 2;
		        i++;
		      }
		    }

		    //return rc.text.length + add_newlines;

		    //We need to substract the no. of lines
		    return rc.text.length - add_newlines; 
		  }  
		  return 0; 
		}

		$scope.selectbyRating = function(rating){
			$scope.libsView = [];
			if (rating != null){
				angular.forEach($scope.allLibs,function(v,k){
					if (v.rating == rating.value){
						$scope.libsView.push(v);
					}
				});
			}else{
				$scope.libsView = $scope.allLibs;
			}
		}

	});
});