<?php
declare (strict_types = 1);

# Copyright (c) 2015, phpfmt and its authors
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
#
# 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
#
# 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
#
# 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

namespace {
	$concurrent = function_exists('pcntl_fork');
	if ($concurrent) {
		require 'vendor/ccirello/csp/csp.php';
	}
	require 'Core/Cacher.php';
	$enableCache = false;
	if (class_exists('SQLite3')) {
		$enableCache = true;
		require 'Core/Cache.php';
	}
	require 'Core/CacheDummy.php';

	require 'version.php';
	require 'helpers.php';
	require 'selfupdate.php';

	require 'Core/constants.php';
	require 'Core/FormatterPass.php';
	require 'Additionals/AdditionalPass.php';
	require 'Core/BaseCodeFormatter.php';
	if ('1' === getenv('FMTDEBUG') || 'step' === getenv('FMTDEBUG')) {
		require 'Core/CodeFormatter_debug.php';
	} elseif ('profile' === getenv('FMTDEBUG')) {
		require 'Core/CodeFormatter_profile.php';
	} else {
		require 'Core/CodeFormatter.php';
	}

	require 'Core/AddMissingCurlyBraces.php';
	require 'Core/AutoImport.php';
	require 'Core/DoubleToSingleQuote.php';
	require 'Core/EliminateDuplicatedEmptyLines.php';
	require 'Core/ExtraCommaInArray.php';
	require 'Core/IndentTernaryConditions.php';
	require 'Core/MergeCurlyCloseAndDoWhile.php';
	require 'Core/MergeDoubleArrowAndArray.php';
	require 'Core/MergeElseIf.php';
	require 'Core/MergeParenCloseWithCurlyOpen.php';
	require 'Core/NormalizeIsNotEquals.php';
	require 'Core/NormalizeLnAndLtrimLines.php';
	require 'Core/PSR2ModifierVisibilityStaticOrder.php';
	require 'Core/Reindent.php';
	require 'Core/ReindentColonBlocks.php';
	require 'Core/ReindentComments.php';
	require 'Core/ReindentEqual.php';
	require 'Core/ReindentObjOps.php';
	require 'Core/RemoveBOMMark.php';
	require 'Core/ResizeSpaces.php';
	require 'Core/ReturnNull.php';
	require 'Core/RTrim.php';
	require 'Core/ShortArray.php';
	require 'Core/SplitCurlyCloseAndTokens.php';
	require 'Core/StripExtraCommaInList.php';
	require 'Core/SurrogateToken.php';
	require 'Core/TwoCommandsInSameLine.php';

	require 'Additionals/AddMissingParentheses.php';
	require 'Additionals/AlignConstVisibilityEquals.php';
	require 'Additionals/AlignDoubleArrow.php';
	require 'Additionals/AlignDoubleSlashComments.php';
	require 'Additionals/AlignEquals.php';
	require 'Additionals/AlignGroupDoubleArrow.php';
	require 'Additionals/AlignPHPCode.php';
	require 'Additionals/AlignTypehint.php';
	require 'Additionals/ConvertOpenTagWithEcho.php';
	require 'Additionals/DocBlockToComment.php';
	require 'Additionals/EncapsulateNamespaces.php';
	require 'Additionals/LeftWordWrap.php';
	require 'Additionals/MergeNamespaceWithOpenTag.php';
	require 'Additionals/NewLineBeforeReturn.php';
	require 'Additionals/NoSpaceAfterPHPDocBlocks.php';
	require 'Additionals/OrganizeClass.php';
	require 'Additionals/OrderAndRemoveUseClauses.php';
	require 'Additionals/OnlyOrderUseClauses.php';
	require 'Additionals/PHPDocTypesToFunctionTypehint.php';
	require 'Additionals/PrettyPrintDocBlocks.php';
	require 'Additionals/ReindentAndAlignObjOps.php';
	require 'Additionals/ReindentSwitchBlocks.php';
	require 'Additionals/RemoveIncludeParentheses.php';
	require 'Additionals/RemoveSemicolonAfterCurly.php';
	require 'Additionals/RemoveUseLeadingSlash.php';
	require 'Additionals/RestoreComments.php';
	require 'Additionals/SmartLnAfterCurlyOpen.php';
	require 'Additionals/SortUseNameSpace.php';
	require 'Additionals/SpaceAroundControlStructures.php';
	require 'Additionals/SpaceBetweenMethods.php';
	require 'Additionals/StripNewlineAfterClassOpen.php';
	require 'Additionals/StripNewlineAfterCurlyOpen.php';
	require 'Additionals/StripNewlineWithinClassBody.php';
	require 'Additionals/StripSpaceWithinControlStructures.php';
	require 'Additionals/TrimSpaceBeforeSemicolon.php';
	require 'Additionals/WordWrap.php';
	require 'Additionals/YodaComparisons.php';

	if (!isset($inPhar)) {
		$inPhar = false;
	}
	if (!isset($testEnv)) {
		require 'cli-core.php';
	}
}