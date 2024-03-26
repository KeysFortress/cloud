using System.Security.Cryptography.X509Certificates;
using System.Text;
using Domain;
using Infrastructure;
using NSec.Cryptography;

namespace Application;

public class AuthenticationService : IAuthenticationService
{
    private readonly List<StoredChallenge> _storedChallenges = new();
    private readonly Timer _timer;


    public AuthenticationService()
    {
        _timer = new Timer(OnTokensExpire, null, TimeSpan.Zero, TimeSpan.FromSeconds(10));

    }

    private void OnTokensExpire(object? state)
    {
        _storedChallenges.RemoveAll(x => x.ExpiresAt < DateTime.UtcNow);
    }

    public AuthenticationResponse InitLogin(string email)
    {
        var userId = Guid.NewGuid(); // Simulate user ID
        var id = Guid.NewGuid();
        var uniqueMessage = GenerateRandomCode();
        var storedMessage = new StoredChallenge
        {
            Id = id.ToString(),
            Challenge = uniqueMessage,
            ExpiresAt = DateTime.UtcNow.AddMinutes(10)
        };
        _storedChallenges.Add(storedMessage);

        return new AuthenticationResponse
        {
            Challenge = uniqueMessage,
            RpId = "http://localhost:5263/",
            Timeout = 60000,
            UserVerification = "preffered"
        };
    }

    public StoredChallenge? VerifyChallenge(byte[] signature, string challenge, byte[] publicKeyBytes)
    {
        var signatureBytes = signature;
        var algorithm = new Ed25519();
        var publicKey = NSec.Cryptography.PublicKey.Import(algorithm, publicKeyBytes, KeyBlobFormat.RawPublicKey);
        var last = _storedChallenges.LastOrDefault(x => x.Id == challenge);
        if (last == null) return null;

        var messageBytes = Encoding.UTF8.GetBytes(last.Challenge);
        var isSignatureValid = algorithm.Verify(publicKey, messageBytes, signatureBytes);


        return isSignatureValid ? last : null;
    }

    public Task<bool> LogOut()
    {
        throw new NotImplementedException();
    }

    private static string GenerateRandomCode()
    {
        const string chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
        var random = new Random();

        var code = new StringBuilder(8);
        for (var i = 0; i < 8; i++)
        {
            code.Append(chars[random.Next(chars.Length)]);
        }

        return code.ToString();
    }

    public bool Login(Guid userId)
    {
        return true;
    }
}
