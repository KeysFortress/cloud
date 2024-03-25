using Microsoft.AspNetCore.Components;

namespace Blazor.Components.Panels;

public partial class WelcomeScreen
{
    String Email { get; set; } = "";

    [Parameter] public EventCallback<string> OnEmailEntered { get; set; }

    void OnNextPressed()
    {
        if (string.IsNullOrEmpty(Email)) return;

        OnEmailEntered.InvokeAsync(Email);
    }
}
