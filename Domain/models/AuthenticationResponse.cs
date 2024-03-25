namespace Domain;

public class AuthenticationResponse
{
    public required string Challenge { get; set; }
    public required int Timeout { get; set; }
    public required string RpId { get; set; }
    public required string UserVerification { get; set; }
}
