using Domain;

namespace Infrastructure;

public interface IAuthenticationService
{
    public Task<AuthenticationResponse> InitLogin();
    public Task<bool> Login(byte[] signature, string challenge, byte[] publicKey);
    public Task<bool> LogOut();
}
