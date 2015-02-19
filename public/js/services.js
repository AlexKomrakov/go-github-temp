(function () {
    angular
        .module('gohub')
        .service('reposService', reposService);

    function reposService($resource) {
        var reposService = this;
        var repositories = { data: [] };
        reposService.api = $resource("/repos", {},  {
            setHook: {method:'GET', url: "/repos/:user/:repo/hook", params:{user: "@user", repo: "@repo"}}
        });
        reposService.loadRepositories = function() {
            reposService.api.query(function(result){
                repositories.data = result;
            })
        };
        reposService.getRepositories = function () {
            reposService.loadRepositories();
            return repositories;
        };
        reposService.addRepository = function(repo) {
            reposService.api.save(repo, function(){
                reposService.loadRepositories();
            })
        };
        reposService.setHook = function(user, repo) {
            var result = { data: null };
            reposService.api.setHook({user: user, repo: repo}).$promise.then(function(res){
                result.data = res;
            });
            return result
        };
    }
})();