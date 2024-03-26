using Domain;

namespace Infrastructure;

public interface IAuthenticationService
{
    public AuthenticationResponse InitLogin(string email);
    public StoredChallenge? VerifyChallenge(byte[] signature, string challenge, byte[] publicKeyBytes);
    public bool Login(Guid userId);
    public Task<bool> LogOut();
}
