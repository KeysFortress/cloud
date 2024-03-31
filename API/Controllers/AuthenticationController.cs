using Infrastructure;
using Microsoft.AspNetCore.Mvc;

namespace API;


[Route("/API/V1/[controller]")]
public class AuthenticationController : BaseApiController
{
    readonly IAuthenticationService _authentinicationService;
    public AuthenticationController(IAuthenticationService authenticationService)
    {
        _authentinicationService = authenticationService;
    }

    [HttpGet("Init-Login/{email}")]
    public IActionResult InitLogin(string email)
    {
        var challenge = _authentinicationService.InitLogin(email);
        return Ok(challenge);
    }

    [HttpPost("Finish")]
    public IActionResult FinishAuthentinication([FromBody] Guid id)
    {
        var isVerified = _authentinicationService.IsChallengeVerified(id);
        return Ok(isVerified);
    }
}
