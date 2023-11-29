import {AsyncPipe, NgClass, NgForOf} from '@angular/common';
import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';

import {Subject} from 'rxjs';

import {BingoCardDto} from '../../core/api/model/bingoCardDto.model';
import {BingoApiService} from '../../core/api/service/bingo-api.service';
import {UserCountPipe} from '../../core/pipes/user-count.pipe';
import {RefreshService} from '../../core/services/refresh.service';

@Component({
    standalone: true,
    selector: 'app-private-bingo-board',
    styleUrls: ['private-bingo-board.page.component.scss'],
    imports: [AsyncPipe, NgForOf, UserCountPipe, NgClass],
    templateUrl: 'private-bingo-board.page.component.html'
})
export class PrivateBingoBoardPageComponent implements OnInit {
    private readonly boardUuid: string;
    public bingoCardsSubject = new Subject<BingoCardDto[]>();

    constructor(
        private route: ActivatedRoute,
        private refreshService: RefreshService,
        private bingoApiService: BingoApiService
    ) {
        this.boardUuid = this.route.snapshot.params['uuid'];
    }

    ngOnInit(): void {
        this.refreshService.onRefreshBingoCards().subscribe(() => this.fetchBingoCards());
        this.fetchBingoCards();
    }

    private fetchBingoCards() {
        console.log('Fetch');
        this.bingoApiService.getBingoBoard(this.boardUuid).subscribe(board => this.bingoCardsSubject.next(board.bingo_cards));
    }

    clickBingoCard(id: string): void {
        this.bingoApiService.clickBingoCard(id).subscribe(() => this.fetchBingoCards());
    }
}
