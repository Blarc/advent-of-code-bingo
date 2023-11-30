import {AsyncPipe, NgClass, NgForOf} from '@angular/common';
import {Component, Input, inject} from '@angular/core';

import {Subject} from 'rxjs';

import {BingoCardDto} from '../../api/model/bingoCardDto.model';
import {BingoApiService} from '../../api/service/bingo-api.service';
import {UserCountPipe} from '../../pipes/user-count.pipe';
import {RefreshService} from '../../services/refresh.service';

@Component({
    standalone: true,
    selector: 'app-bingo-grid',
    styleUrls: ['bingo-grid.component.scss'],
    imports: [AsyncPipe, NgForOf, NgClass, UserCountPipe],
    templateUrl: 'bingo-grid.component.html'
})
export class BingoGridComponent {
    private refreshService = inject(RefreshService);
    private bingoApiService = inject(BingoApiService);

    @Input() bingoCards = new Subject<BingoCardDto[]>();

    public clickBingoCard(id: string): void {
        this.bingoApiService.clickBingoCard(id).subscribe(() => this.refreshService.shouldRefreshBingoCards());
    }
}
