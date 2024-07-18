const checkbox_ids = ['AnthropologyEdit', 'ArchaeologyEdit', 'DrivingEdit', 'LibraryResearchEdit', 'AccountingEdit', 'CharmeEdit', 'CthulhuMythosEdit', 'IntimidateEdit', 
    'ElectricRepairsEdit', 'FirstAidEdit', 'FinancialsEdit', 'HistoryEdit', 'ListeningEdit', 'ConcealingEdit', 'ClimbingEdit', 'MechanicalRepairsEdit', 'MedicineEdit',
    'NaturalHistoryEdit', 'OccultismEdit', 'OrientationEdit', 'PsychoAnalysisEdit', 'PsychologyEdit', 'LawEdit', 'HorseridingEdit', 'LocksEdit', 'HeavyMachineryEdit', 
    'SwimmingEdit', 'JumpingEdit', 'TrackingEdit', 'PersuasionEdit', 'ConvincingEdit', 'StealthEdit', 'DetectingSecretsEdit', 'DisguisingEdit', 'ThrowingEdit', 'ValuationEdit'];

for (const id of checkbox_ids){
    var el = document.getElementById(id);
    if (el) {
        el.addEventListener('click', function() {
            newId = id.substring(0, id.length - 4);
            elem = document.getElementById(newId)
            elem.readOnly = !this.checked;
        })
    }
}
