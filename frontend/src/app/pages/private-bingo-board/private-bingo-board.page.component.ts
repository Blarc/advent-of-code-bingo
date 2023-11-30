import {AsyncPipe, NgClass, NgForOf, NgIf} from '@angular/common';
import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';

import {Subject} from 'rxjs';

import {BingoBoardDto} from '../../core/api/model/bingoBoardDto.model';
import {BingoCardDto} from '../../core/api/model/bingoCardDto.model';
import {BingoApiService} from '../../core/api/service/bingo-api.service';
import {BingoGridComponent} from '../../core/components/bingo-grid/bingo-grid.component';
import {UserCountPipe} from '../../core/pipes/user-count.pipe';
import {RefreshService} from '../../core/services/refresh.service';

@Component({
    standalone: true,
    selector: 'app-private-bingo-board',
    styleUrls: ['private-bingo-board.page.component.scss'],
    imports: [AsyncPipe, NgForOf, UserCountPipe, NgClass, NgIf, BingoGridComponent],
    templateUrl: 'private-bingo-board.page.component.html'
})
export class PrivateBingoBoardPageComponent implements OnInit {
    private readonly boardUuid: string;
    public bingoBoardSubject = new Subject<BingoBoardDto>();
    public bingoCardsSubject = new Subject<BingoCardDto[]>();

    constructor(
        private route: ActivatedRoute,
        private refreshService: RefreshService,
        private bingoApiService: BingoApiService
    ) {
        this.boardUuid = this.route.snapshot.params['uuid'];
    }

    ngOnInit(): void {
        this.fetchBingoCards();
        this.refreshService.onRefreshBingoCards().subscribe(() => this.fetchBingoCards());
    }

    private fetchBingoCards() {
        this.bingoApiService.getBingoBoard(this.boardUuid).subscribe(board => {
            this.bingoBoardSubject.next(board);
            this.bingoCardsSubject.next(board.bingo_cards);
        });
    }
}
