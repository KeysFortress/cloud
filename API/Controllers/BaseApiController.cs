using System.Security.Claims;
using Microsoft.AspNetCore.Mvc;

namespace API;

public class BaseApiController : ControllerBase
{

    public BaseApiController()
    {
    }

    protected Guid GetRequestOwner()
    {
        var user = User as ClaimsPrincipal;
        var claim = user?.FindFirst("uuid")?.Value;
        return claim == null ? Guid.Empty : Guid.Parse(claim);
    }
}
