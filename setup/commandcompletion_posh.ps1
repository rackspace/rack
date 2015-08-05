function global:TabExpansion2 {
	[CmdletBinding(DefaultParameterSetName = 'ScriptInputSet')]
	Param(
    		[Parameter(ParameterSetName = 'ScriptInputSet', Mandatory = $true, Position = 0)]
    		[string] $inputScript,

    		[Parameter(ParameterSetName = 'ScriptInputSet', Mandatory = $true, Position = 1)]
    		[int] $cursorColumn,

    		[Parameter(ParameterSetName = 'AstInputSet', Mandatory = $true, Position = 0)]
    		[System.Management.Automation.Language.Ast] $ast,

    		[Parameter(ParameterSetName = 'AstInputSet', Mandatory = $true, Position = 1)]
    		[System.Management.Automation.Language.Token[]] $tokens,

    		[Parameter(ParameterSetName = 'AstInputSet', Mandatory = $true, Position = 2)]
    		[System.Management.Automation.Language.IScriptPosition] $positionOfCursor,

    		[Parameter(ParameterSetName = 'ScriptInputSet', Position = 2)]
    		[Parameter(ParameterSetName = 'AstInputSet', Position = 3)]
    		[Hashtable] $options = $null
	)

	End {
    $result = $null

    if ($psCmdlet.ParameterSetName -eq 'ScriptInputSet') {
      $result = [System.Management.Automation.CommandCompletion]::CompleteInput(
        <#inputScript#>  $inputScript,
        <#cursorColumn#> $cursorColumn,
        <#options#>      $options)
    }
    else{
      $result = [System.Management.Automation.CommandCompletion]::CompleteInput(
        <#ast#>              $ast,
        <#tokens#>           $tokens,
        <#positionOfCursor#> $positionOfCursor,
        <#options#>          $options)
    }


    if ($result.CompletionMatches.Count -eq 0){
			if ($psCmdlet.ParameterSetName -eq 'ScriptInputSet') {
        $ast = [System.Management.Automation.Language.Parser]::ParseInput($inputScript, [ref]$tokens, [ref]$null)
      }
      $text = $ast.Extent.Text
    	if($text -match '^*rack.exe*') {
        $cmd1 = $text -split '\s+'
        $end = $cmd1.count - 2
        $cmd2 = $cmd1[0..$end]
        $cmd3 = $cmd2 -join ' '
        $suggestions = Invoke-Expression "$cmd3 --generate-bash-completion"
        ForEach($suggestion in $suggestions) {
          if($suggestion -match $cmd1[$end + 1]) {
            $suggestionObject = New-Object System.Management.Automation.CompletionResult ($suggestion, $suggestion, "Text", $suggestion)
				      $result.CompletionMatches.Add($suggestionObject)
          }
        }
    	}
		}

		return	$result
	}
}
