using Domain;
using Infrastructure;

namespace Application;

public class AuthenticationService : IAuthenticationService
{
    public async Task<AuthenticationResponse> InitLogin()
    {
        return new AuthenticationResponse
        {
            Challenge = "dwadwad wad wad ",
            RpId = "http://localhost:5263/",
            Timeout = 60000,
            UserVerification = "preffered"
        };
    }

    public async Task<bool> Login(byte[] signature, string challenge, byte[] publicKey)
    {
        return true;
    }

    public Task<bool> LogOut()
    {
        throw new NotImplementedException();
    }
}
