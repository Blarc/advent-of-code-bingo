import {AsyncPipe, NgClass, NgForOf, NgIf} from '@angular/common';
import {Component, OnInit, inject} from '@angular/core';

import {Subject} from 'rxjs';

import {BingoCardDto} from '../../core/api/model/bingoCardDto.model';
import {BingoApiService} from '../../core/api/service/bingo-api.service';
import {UserCountPipe} from '../../core/pipes/user-count.pipe';
import {RefreshService} from '../../core/services/refresh.service';

@Component({
    standalone: true,
    selector: 'app-bingo-page',
    templateUrl: 'bingo.page.component.html',
    imports: [NgForOf, NgIf, NgClass, AsyncPipe, UserCountPipe],
    styleUrls: ['bingo.page.component.scss']
})
export class BingoPageComponent implements OnInit {
    private refreshService = inject(RefreshService);
    private bingoApiService = inject(BingoApiService);

    public bingoCardsSubject = new Subject<BingoCardDto[]>();

    ngOnInit(): void {
        this.refreshService.onRefreshBingoCards().subscribe(() => this.fetchBingoCards());
        this.fetchBingoCards();
    }

    private fetchBingoCards() {
        this.bingoApiService.getAllBingoCards().subscribe(cards => this.bingoCardsSubject.next(cards));
    }

    clickBingoCard(id: string): void {
        this.bingoApiService.clickBingoCard(id).subscribe(() => this.fetchBingoCards());
    }
}
