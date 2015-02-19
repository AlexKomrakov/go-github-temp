(function () {

    angular
        .module('gohub')
        .controller('indexController', indexController)
        .controller('repoController', repoController);

    function indexController($scope, reposService) {
        $scope.repositories = reposService.getRepositories();
        $scope.addRepository = function(repo) {
            reposService.addRepository(repo)
        }
    }

    function repoController($scope, reposService) {
        $scope.setHook = function (user, repo) {
            $scope.okHook = reposService.setHook(user, repo);
        }
    }

})();