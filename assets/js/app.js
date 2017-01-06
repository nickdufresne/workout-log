var app = angular.module("Wolo", ["ui.router", "ngResource"]);

app.config(function($stateProvider, $urlRouterProvider) {
  //
  // For any unmatched url, redirect to /state1
  $urlRouterProvider.otherwise("/");
  //
  // Now set up the states
  $stateProvider
    .state('home', {
      url: "/",
      templateUrl: "/views/index.html",
      controller: 'WorkoutIndexController',
    })
    .state('settings', {
      url: "/settings",
      templateUrl: "/views/settings.html",
      controller: 'SettingsController',
    })
    .state('newWorkout', {
      url: "/workouts/new",
      templateUrl: "/views/workouts/new.html",
      controller: 'NewWorkoutController'
    })
    .state('viewWorkout', {
      url: "/workouts/:workoutID",
      templateUrl: "/views/workouts/view.html",
      controller: 'ViewWorkoutController'
    })
    .state('editWorkout', {
      url: "/workouts/:workoutID/edit",
      templateUrl: "/views/workouts/edit.html",
      controller: 'EditWorkoutController'
    });
});

app.factory("$workout", function($resource){
  return $resource("/workouts/:workoutID", {workoutID: "@workoutID"})
});

app.controller("WorkoutIndexController", function($scope){
  $scope.message = "Wolo";
});

app.controller("NewWorkoutController", function($scope){
  $scope.message = "Wolo";
});

app.controller("ViewWorkoutController", function($scope){
  $scope.message = "Wolo";
});

app.controller("EditWorkoutController", function($scope){
  $scope.message = "Wolo";
});

app.controller("SettingsController", function($scope, $resource, $state){

  var Settings = $resource('/settings', {});
  
  $scope.settings = Settings.get({});

  $scope.saveSettings = function(){
    Settings.save(
      {name: $scope.settings.name, email: $scope.settings.email}, 
      function(){
        $state.go('home');
      }
    );
  };
});