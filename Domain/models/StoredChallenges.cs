namespace Domain;

public class StoredChallenge
{
    public Guid Id { get; set; }
    public string Challenge { get; set; }
    public DateTime ExpiresAt { get; set; }
    public Guid userId { get; set; }
    public bool Verified { get; set; }
}
