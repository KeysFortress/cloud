using Infrastructure;
using Microsoft.AspNetCore.Components;
using Microsoft.JSInterop;

namespace Blazor.Components.Panels;

public partial class LoginPanel
{
    [Inject]
    IJSRuntime JSRuntime { get; set; }
    [Inject]
    NavigationManager NavigationManager { get; set; }
    [Inject]
    IAuthenticationService AuthenticationService { get; set; }

    [Parameter] public string Email { get; set; }
    protected override async Task OnAfterRenderAsync(bool firstRender)
    {
        if (firstRender)
        {
            try
            {
                var code = await AuthenticationService.InitLogin();
                Console.WriteLine(code);
                await JSRuntime.InvokeVoidAsync("onRenderComplete", null);

            }
            catch (Exception ex)
            {

            }
        }
    }

    void OnLoginPressed()
    {
        NavigationManager.NavigateTo("dashboard");
    }
}
