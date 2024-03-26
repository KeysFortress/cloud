using Domain;

namespace Infrastructure;

public interface IAuthenticationService
{
    public AuthenticationResponse InitLogin(string email);
    public StoredChallenge? VerifyChallenge(byte[] signature, Guid id, byte[] publicKeyBytes);
    bool IsChallengeVerified(Guid id);
    public bool Login(Guid userId);
    public Task<bool> LogOut();
}
